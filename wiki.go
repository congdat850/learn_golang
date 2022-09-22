// example 1

// package main

// import (
// 	"fmt"
// 	"os"
// )

// type Page struct {
// 	Title string
// 	Body  []byte
// }

// func (p *Page) save() error { // create object function
// 	filename := p.Title + ".txt"
// 	return os.WriteFile(filename, p.Body, 0600)
// }

// func loadPage(title string) (*Page, error) {
// 	filename := title + ".txt"
// 	body, err := os.ReadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Page{Title: title, Body: body}, nil
// }

// func main() {
// 	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
// 	p1.save()
// 	p2, _ := loadPage("TestPage")
// 	fmt.Println(string(p2.Body))
// }

// example 2
// go:build ignore

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// func handler(w http.ResponseWriter, r *http.Request) { // handl request and take response
// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

// func main() {
// 	http.HandleFunc("/", handler)
// 	log.Fatal(http.ListenAndServe(":8080", nil)) // log error
// }

// example 3

// package main

// import (
// 	// "fmt"

// 	"html/template"
// 	"log"
// 	"net/http"
// 	"os"
// )

// type Page struct {
// 	Title string
// 	Body  []byte
// }

// func loadPage(title string) (*Page, error) {
// 	filename := title + ".txt"
// 	body, err := os.ReadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Page{Title: title, Body: body}, nil
// }

// func (p *Page) save() error { // create object function
// 	filename := p.Title + ".txt"
// 	return os.WriteFile(filename, p.Body, 0600)
// }

// func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
// 	// t, _ := template.ParseFiles(tmpl + ".html")
// 	// t.Execute(w, p)

// 	t, err := template.ParseFiles(tmpl + ".html")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	err = t.Execute(w, p)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func viewHandler(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/view/"):]
// 	p, err := loadPage(title)
// 	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)

// 	// t, _ := template.ParseFiles("view.html")
// 	// t.Execute(w, p)

// 	if err != nil {
// 		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
// 		return
// 	}

// 	renderTemplate(w, "view", p)
// }

// func editHandler(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/edit/"):]
// 	p, err := loadPage(title)
// 	if err != nil {
// 		p = &Page{Title: title}
// 	}
// 	// fmt.Fprintf(w, "<h1>Editing %s</h1>"+
// 	// 	"<form action=\"/save/%s\" method=\"POST\">"+
// 	// 	"<textarea name=\"body\">%s</textarea><br>"+
// 	// 	"<input type=\"submit\" value=\"Save\">"+
// 	// 	"</form>",
// 	// 	p.Title, p.Title, p.Body)

// 	// t, _ := template.ParseFiles("edit.html")
// 	// t.Execute(w, p)

// 	renderTemplate(w, "edit", p)
// }

// func saveHandler(w http.ResponseWriter, r *http.Request) {
// 	// title := r.URL.Path[len("/save/"):]
// 	// body := r.FormValue("body")
// 	// p := &Page{Title: title, Body: []byte(body)}
// 	// p.save()
// 	// http.Redirect(w, r, "/view/"+title, http.StatusFound)

// 	title := r.URL.Path[len("/save/"):]
// 	body := r.FormValue("body")
// 	p := &Page{Title: title, Body: []byte(body)}
// 	err := p.save()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	http.Redirect(w, r, "/view/"+title, http.StatusFound)
// }
// func main() {
// 	http.HandleFunc("/view/", viewHandler)
// 	http.HandleFunc("/edit/", editHandler)
// 	http.HandleFunc("/save/", saveHandler)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// exmaple 4

package main

import (
	// "fmt"

	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func (p *Page) save() error { // create object function
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	// t, _ := template.ParseFiles(tmpl + ".html")
	// t.Execute(w, p)

	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
