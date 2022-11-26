package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/add-project", project).Methods("GET")
	route.HandleFunc("/article", article).Methods("GET")

	var port string = "5000"
	fmt.Print("Server running on port " + port)
	http.ListenAndServe("localhost:"+port, route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8") //fungsinya?
	tmpt, err := template.ParseFiles("index.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8") //fungsinya?
	tmpt, err := template.ParseFiles("contact.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8") //fungsinya?
	tmpt, err := template.ParseFiles("myProject.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

func article(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8") //fungsinya?
	tmpt, err := template.ParseFiles("article.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}
