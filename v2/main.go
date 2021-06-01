package main

import (
	"log"
	"net/http"
	"text/template"
	"time"
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

func init() {
	tpl = template.Must(template.ParseGlob("./*.html"))
}
func show(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(res, "front_template.html", posts)
	if err != nil {
		log.Fatal(err)
	}
}

func send(res http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	note := req.FormValue("note")
	if name != "" && note != "" {
		posts = append(posts, post{name, note, time.Now()})
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(res, "front_template.html", posts)
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
	err := tpl.ExecuteTemplate(res, "front_template.html", data)
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

func main() {
	http.Handle("/names", http.HandlerFunc(names))
	http.Handle("/", http.HandlerFunc(show))
	http.Handle("/show", http.HandlerFunc(show))
	http.Handle("/post", http.HandlerFunc(send))
	http.Handle("/search", http.HandlerFunc(search))
	http.ListenAndServe(":8080", nil)

}
