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

func front(res http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	note := req.FormValue("note")
	if name != "" && note != "" {
		posts = append(posts, post{name, note, time.Now()})
	}
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
		data = posts
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.ExecuteTemplate(res, "front_template.html", data)
	if err != nil {
		log.Fatal(err)
	}
	post_byname = nil
}

func main() {

	http.Handle("/", http.HandlerFunc(front))
	http.ListenAndServe(":8080", nil)

}
