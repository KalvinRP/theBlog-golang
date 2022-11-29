package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/asset").Handler(http.StripPrefix("/asset", http.FileServer(http.Dir("./asset"))))

	route.HandleFunc("/article/{ID}", article).Methods("GET")
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/add-project", project).Methods("GET")
	route.HandleFunc("/add-project", addprojects).Methods("POST")
	route.HandleFunc("/delete/{index}", delete).Methods("GET")

	var port string = "5000"
	fmt.Print("Server running on port " + port)
	http.ListenAndServe("localhost:"+port, route)
}

type Prj struct {
	PrjName  string
	Duration int
	Desc     string
	Tech     []string
	// img      string
}

var Addprj = []Prj{
	{
		PrjName:  "Placeholder",
		Duration: 22,
		Desc:     "Lorem ipsum",
		// Tech:     techno,
		// img:      fileLocation,
	},
}

func addprojects(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	// err := r.ParseMultipartForm(1024)
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

	Duration := t2.Sub(t1)
	duratext := int(Duration.Hours() / 24)

	Desc := r.Form.Get("desc")

	techno := r.Form["tech"]

	// image, imgname, err := r.FormFile("image")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// defer image.Close()
	// dir, err := os.Getwd()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// filename := imgname.Filename
	// fileLocation := filepath.Join(dir, "submittedImage", filename)
	// targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// defer targetFile.Close()
	// if _, err := io.Copy(targetFile, image); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//fileupload
	var newprj = Prj{
		PrjName:  pname,
		Duration: duratext,
		Desc:     Desc,
		Tech:     techno,
		// img:      fileLocation,
	}

	// fmt.Println(newprj) //Result: Working

	Addprj = append(Addprj, newprj)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("index.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	project := map[string]interface{}{
		"Project": Addprj,
	}

	// fmt.Println(project) //Result: Working

	tmpt.Execute(w, project)
}

func delete(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	fmt.Println(index)

	Addprj = append(Addprj[:index], Addprj[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
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

func article(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("article.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	ID, _ := strconv.Atoi(mux.Vars(r)["ID"])

	var Detail = Prj{}

	for index, data := range Addprj {
		if index == ID {
			Detail = Prj{
				PrjName:  data.PrjName,
				Duration: data.Duration,
				Desc:     data.Desc,
				Tech:     data.Tech,
				// img:   data.img,
			}
		}
	}

	article := map[string]interface{}{
		"Article": Detail,
	}

	tmpt.Execute(w, article)
}
