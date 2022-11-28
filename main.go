package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/asset").Handler(http.StripPrefix("/asset", http.FileServer(http.Dir("./asset"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/add-project", project).Methods("GET")
	route.HandleFunc("/add-project", addprojects).Methods("POST")
	route.HandleFunc("/article", article).Methods("GET")

	var port string = "5000"
	fmt.Print("Server running on port " + port)
	http.ListenAndServe("localhost:"+port, route)
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

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("index.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	addprj := map[string]interface{}{
		"Project": addprj,
	}

	fmt.Println(addprj)

	tmpt.Execute(w, addprj)
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
	prjName  string
	duration float64
	desc     string
	tech     []string
	img      string
}

var addprj = []prj{{}}

func addprojects(w http.ResponseWriter, r *http.Request) {
	// err := r.ParseForm()

	err := r.ParseMultipartForm(1024)
	if err != nil {
		log.Fatal(err)
	}

	pname := r.Form.Get("prjname")

	const format string = "2006-1-2 15:04:05"

	var sdate string = r.Form.Get("sdate")
	start := sdate + " 00:00:00"       //ditambah jam supaya bisa jalan
	t1, _ := time.Parse(format, start) //reflect.TypeOf(t1) = time.Time

	var edate string = r.Form.Get("edate")
	end := edate + " 00:00:00"
	t2, _ := time.Parse(format, end) //reflect.TypeOf(t2) = time.Time

	duration := t2.Sub(t1)
	duratext := duration.Hours() / 24

	desc := r.Form.Get("desc")

	techno := r.Form["tech"]

	image, imgname, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer image.Close()
	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := imgname.Filename
	fileLocation := filepath.Join(dir, "submittedImage", filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()
	if _, err := io.Copy(targetFile, image); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fileupload
	var newprj = prj{
		prjName:  pname,
		duration: duratext,
		desc:     desc,
		tech:     techno,
		img:      fileLocation,
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
