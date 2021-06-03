package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	uuid "github.com/satori/go.uuid"
)

var tpl *template.Template

type post struct {
	Name string
	Post string
	Date time.Time
}

var posts = make([]post, 0)
var post_byname = make([]post, 0)
var data = make([]post, 0)
var users = make(map[string]string)
var sessions = make(map[string]string)

func init() {
	tpl = template.Must(template.ParseGlob("./*.html"))
}
func show(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(res, "noteable.html", posts)
	if err != nil {
		log.Fatal(err)
	}
}

func send(res http.ResponseWriter, req *http.Request) {
	note := req.FormValue("note")
	cookie, err := req.Cookie("session")
	if err != nil {
		log.Fatal(err)
	}
	id := cookie.Value
	name := sessions[id]
	if note != "" {
		posts = append(posts, post{name, note, time.Now()})
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tpl.ExecuteTemplate(res, "noteable.html", posts)
	if err != nil {
		log.Fatal(err)
	}
}

func search(res http.ResponseWriter, req *http.Request) {
	search_name := req.FormValue("search_name")
	if search_name != "" {
		for _, note := range posts {
			if note.Name == search_name {
				post_byname = append(post_byname, note)
			}
		}
	}
	if len(post_byname) > 0 {
		data = post_byname
	} else {
		data = []post{{"not found!", "", time.Now()}}
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(res, "noteable.html", data)
	if err != nil {
		log.Fatal(err)
	}
	post_byname = nil
}

func names(res http.ResponseWriter, req *http.Request) {
	var names []string
	for _, note := range posts {
		name := note.Name
		if !find(names, name) {
			names = append(names, name)
		}

	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(res, "names.html", names)
	if err != nil {
		log.Fatal(err)
	}

}

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func registery(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(res, "registery.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func register(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	name := req.FormValue("name")
	password := req.FormValue("password")

	if _, ok := users[name]; ok {
		err := tpl.ExecuteTemplate(res, "registery.html", name+" already exists")
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	users[name] = password
	err := tpl.ExecuteTemplate(res, "front.html", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(users)
}

func loginery(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(res, "login.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func login(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	name := req.FormValue("name")
	password := req.FormValue("password")
	if users[name] == password {
		id := uuid.NewV4()
		cookie := &http.Cookie{
			Name:     "session",
			Value:    id.String(),
			HttpOnly: true,
			Path:     "/",
		}
		sessions[id.String()] = name
		http.SetCookie(res, cookie)

		err := tpl.ExecuteTemplate(res, "noteable.html", nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := tpl.ExecuteTemplate(res, "login.html", "one of the credentials is not correct")
	if err != nil {
		log.Fatal(err)
	}
}

func front(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(res, "front.html", posts)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.Handle("/registery", http.HandlerFunc(registery))
	http.Handle("/register", http.HandlerFunc(register))
	http.Handle("/loginery", http.HandlerFunc(loginery))
	http.Handle("/login", http.HandlerFunc(login))
	http.Handle("/names", http.HandlerFunc(names))
	http.Handle("/", http.HandlerFunc(front))
	http.Handle("/show", http.HandlerFunc(show))
	http.Handle("/send", http.HandlerFunc(send))
	http.Handle("/search", http.HandlerFunc(search))
	http.ListenAndServe(":8080", nil)

}
