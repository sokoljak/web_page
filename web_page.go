package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("mainHandler()")
	p, _ := loadPage("main")
	t, _ := template.ParseFiles("main.html")
	t.Execute(w, p)
}

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("resumeHandler()")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("aboutHandler()")
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/resume", resumeHandler)
	http.HandleFunc("/about", aboutHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
