package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

var ctx = context.Background()

func redisSetSession(name string, uuid string) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, uuid, name, 0).Err()
	if err != nil {
		panic(err)
	}
}

func redisGetSession(uuid string) string {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	name, err := rdb.Get(ctx, uuid).Result()
	if err != nil {
		fmt.Println(err)
	}
	return name

}

var tpl *template.Template

type post struct {
	Name string
	Post string
	Date time.Time
}

func init() {
	tpl = template.Must(template.ParseGlob("./*.html"))
}

func show(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	allposts := getAllPosts()
	err := tpl.ExecuteTemplate(res, "noteable.html", allposts)
	if err != nil {
		log.Fatal(err)
	}
}

func send(res http.ResponseWriter, req *http.Request) {
	note := req.FormValue("note")
	cookie, err := req.Cookie("session")
	var posts = make([]post, 0)
	if err != nil {
		log.Fatal(err)
	}
	id := cookie.Value
	//name := getSession(id)
	name := redisGetSession(id)
	if note != "" {
		setPost(name, note, time.Now())
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
	var data = make([]post, 0)
	if search_name != "" {
		data = getPostByname(search_name)
	}
	if len(data) == 0 {
		data = []post{{"not found!", "", time.Now()}}
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(res, "noteable.html", data)
	if err != nil {
		log.Fatal(err)
	}
}

func names(res http.ResponseWriter, req *http.Request) {
	names := getNames()
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(res, "names.html", names)
	if err != nil {
		log.Fatal(err)
	}

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

	if isUser(name) {
		err := tpl.ExecuteTemplate(res, "registery.html", name+" already exists")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	InsertUser(name, password)
	err := tpl.ExecuteTemplate(res, "front.html", nil)
	if err != nil {
		log.Fatal(err)
	}
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
	if isUserCreds(name, password) {
		id := uuid.NewV4()
		cookie := &http.Cookie{
			Name:     "session",
			Value:    id.String(),
			HttpOnly: true,
			Path:     "/",
		}
		//setSession(id.String(), name)
		redisSetSession(name, id.String())
		http.SetCookie(res, cookie)
		data := getLastPosts(3)
		err := tpl.ExecuteTemplate(res, "noteable.html", data)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err := tpl.ExecuteTemplate(res, "login.html", "one of the credentials is not correct")
	if err != nil {
		log.Fatal(err)
	}
}

func front(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(res, "front.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./style"))
	http.Handle("/style/", http.StripPrefix("/style/", fs))
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
