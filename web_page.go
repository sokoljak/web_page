package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const HtmlExt = ".html"
const TxtExt = ".txt"

var templates = template.Must(template.ParseFiles("main.html", "resume.html"))

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + TxtExt
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + TxtExt
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/main/", http.StatusSeeOther)
}

func mainHandler(w http.ResponseWriter, r *http.Request, title string, p *Page) {
	log.Println("mainHandler()")
	templatesErr := templates.ExecuteTemplate(w, title+HtmlExt, p)
	if templatesErr != nil {
		http.Error(w, strconv.Itoa(http.StatusInternalServerError)+" - "+templatesErr.Error(), http.StatusInternalServerError)
		return
	}
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request, title string, p *Page)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[1 : len(r.URL.Path)-1]
		p, pageLoadErr := loadPage(title)
		if p == nil {
			http.NotFound(w, r)
			return
		} else if pageLoadErr != nil {
			http.Error(w, strconv.Itoa(http.StatusInternalServerError)+" - "+pageLoadErr.Error(), http.StatusInternalServerError)
			return
		} else {
			fn(w, r, title, p)
		}

	}
}

func resumeHandler(w http.ResponseWriter, r *http.Request, title string, p *Page) {
	log.Println("resumeHandler()")
	templatesErr := templates.ExecuteTemplate(w, title+HtmlExt, p)
	if templatesErr != nil {
		http.Error(w, strconv.Itoa(http.StatusInternalServerError)+" - "+templatesErr.Error(), http.StatusInternalServerError)
		return
	}

}

func aboutHandler(w http.ResponseWriter, r *http.Request, title string, p *Page) {
	log.Println("aboutHandler()")
}

func main() {
	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/main/", makeHandler(mainHandler))
	http.HandleFunc("/resume/", makeHandler(resumeHandler))
	http.HandleFunc("/about/", makeHandler(aboutHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
