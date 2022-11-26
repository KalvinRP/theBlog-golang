package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/asset").Handler(http.StripPrefix("/asset", http.FileServer(http.Dir("./asset"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/add-project", project).Methods("GET")
	route.HandleFunc("/add-project", addproject).Methods("POST")
	route.HandleFunc("/article", article).Methods("GET")

	var port string = "5000"
	fmt.Print("Server running on port " + port)
	http.ListenAndServe("localhost:"+port, route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("index.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("contact.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("myProject.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

type prj struct {
	prjName   string
	startDate string
	endDate   string
	desc      string
	tech      string
	img       string
}

var addprj = []prj{
	{
		prjName:   "",
		startDate: "",
		endDate:   "",
		desc:      "",
		tech:      "",
		img:       "",
	},
}

func addproject(w http.ResponseWriter, r *http.Request) {
	eror := r.ParseForm()

	if eror != nil {
		log.Fatal(eror)
	}

	pname := r.PostForm.Get("prjname")
	sdate := r.PostForm.Get("sdate")
	edate := r.PostForm.Get("edate")
	desc := r.PostForm.Get("desc")
	tech := r.PostForm.Get("tech")
	image := r.PostForm.Get("image")

	var newprj = prj{
		prjName:   pname,
		startDate: sdate,
		endDate:   edate,
		desc:      desc,
		tech:      tech,
		img:       image,
	}

	addprj = append(addprj, newprj)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func article(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("article.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}
