package main

import (
	"ascii-art-web/ascii"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var (
	err500 = "500 Internal Server Error"
	err404 = "404 Oops, this page not found..."
	err400 = "400 Bad request"
	err405 = "405 The request method is not allowed"
)

// Page holds font chosen by user, their input, output from art, error
type Page struct {
	Font, Input, Output, Error string
}

var t *template.Template

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, 404)
		return
	}
	// t, _ := template.ParseFiles("index.html")
	t.ExecuteTemplate(w, "index.html", nil)
}

func init() {
	t = template.Must(template.ParseGlob("./templates/*.html"))
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t.ExecuteTemplate(w, "index.html", nil)
	case "POST":
		d := Page{}
		// InputString := r.FormValue("Input")
		r.ParseForm()
		InputS, ok := r.Form["Input"]
		if !ok {
			errorHandler(w, r, 405)
		}
		InputString := strings.Join(InputS, "")
		fn, ok := r.Form["fonts"]
		if !ok {
			errorHandler(w, r, 405)
		}
		InputString = strings.ReplaceAll(InputString, "\r\n", "\n")
		if !isValid(InputString) {
			errorHandler(w, r, 400)
			return
		}
		f := strings.Join(fn, "")
		d.Input = InputString
		d.Font = f

		asciiArt, err := ascii.Art(InputString, f)
		if err != nil {
			d.Error = err.Error()
			errorHandler(w, r, 500)
			return
		}
		d.Output = asciiArt
		// errorHandler(w,r, http.StatusBadRequest)
		t.ExecuteTemplate(w, "index.html", d)
	default:
		errorHandler(w, r, 405)

	}
}

func isValid(s string) bool {
	for _, letter := range s {
		if (letter < 32 || letter > 126) && letter != 10 {
			return false
		}
	}
	return true
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	// t, _ := template.ParseFiles("error.html")
	p := &Page{}
	if status == 404 {
		p = &Page{Error: err404}
	} else if status == 500 {
		p = &Page{Error: err500}
	} else if status == 400 {
		p = &Page{Error: err400}
	} else if status == 405 {
		p = &Page{Error: err405}
	}
	w.WriteHeader(status)
	t.ExecuteTemplate(w, "error.html", p)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/ascii-art", saveHandler)
	fmt.Println("Server start port 3030")
	fmt.Println("http://localhost:3030/")

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":3030", nil)
}
