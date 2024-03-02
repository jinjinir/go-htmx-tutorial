package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	_ "time" // for simulating slow network
)

type Film struct {
	Title    string
	Director string
}

func main() {
	fmt.Println("Hello, World!")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		films := map[string][]Film{
			"Films": {
				{Title: "The Godfather", Director: "Francis Ford Coppola"},
				{Title: "Blade Runner", Director: "Ridley Scott"},
				{Title: "The Thing", Director: "John Carpenter"},
			},
		}
		tmpl.Execute(w, films)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		// time.Sleep(2 * time.Second)  // simulate a slow network
		title := r.PostFormValue("title")
		director := r.PostFormValue("director")

		log.Print("HTMX request received")
		log.Printf("Received title: %s\n", title)
		log.Printf("Received director: %s\n", director)

		// define a template with placeholders for dynamic content
		// TODO: check out how to use template fragments in golang
		tmpl, err := template.New("movie").Parse(
			"<li class='list-group-item bg-primary text-white'>{{ .Title }} - {{ .Director }}</li>")

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// execute the template with the dynamic content
		data := struct {
			Title    string
			Director string
		}{
			Title:    title,
			Director: director,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/add-film/", h2)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
