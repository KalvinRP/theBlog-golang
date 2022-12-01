package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"theblog/connection"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	connection.DatabaseConnect()
	route := mux.NewRouter()

	route.PathPrefix("/asset").Handler(http.StripPrefix("/asset", http.FileServer(http.Dir("./asset"))))
	route.PathPrefix("/images").Handler(http.StripPrefix("/images", http.FileServer(http.Dir("./images"))))

	route.HandleFunc("/article/{id}", article).Methods("GET")
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/add-project", project).Methods("GET")
	route.HandleFunc("/add-project", addprojects).Methods("POST")
	route.HandleFunc("/delete/{id}", delete).Methods("GET")
	route.HandleFunc("/register", register).Methods("GET")
	route.HandleFunc("/register", reginput).Methods("POST")
	route.HandleFunc("/login", login).Methods("GET")
	route.HandleFunc("/login", loginput).Methods("POST")
	route.HandleFunc("/logout", logout).Methods("GET")

	var port string = "5000"
	fmt.Print("Server running on port " + port)
	http.ListenAndServe("localhost:"+port, route)
}

type Prj struct {
	ID         int
	PrjName    string
	Start_date time.Time
	Str_sdate  string
	End_date   time.Time
	Str_edate  string
	Duration   int
	Desc       string
	Tech       []string
	IMG        string
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func addprojects(w http.ResponseWriter, r *http.Request) {
	// err := r.ParseForm()
	err := r.ParseMultipartForm(1024)
	if err != nil {
		log.Fatal(err)
	}

	name := r.Form.Get("prjname")
	sdate := r.Form.Get("sdate")
	edate := r.Form.Get("edate")
	Desc := r.Form.Get("desc")
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
	fileLocation := filepath.Join(dir, "images", filename)
	filenamefit := fileLocation[26:]
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

	// server storage unused due to database connection
	// var newprj = Prj{
	// 	PrjName:  name,
	// 	Duration: duratext,
	// 	Desc:     Desc,
	// 	Tech:     techno,
	// IMG:      fileLocation,
	// }

	// fmt.Println(newprj) //Result: Working

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, description, technologies, Start_date, End_date, image) VALUES ($1, $2, $3, $4, $5, $6)", name, Desc, techno, sdate, edate, filenamefit)
	//belum masukin image
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error on INSERT : " + err.Error()))
		return
	}

	// Addprj = append(Addprj, newprj)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("index.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	newprj, err := connection.Conn.Query(context.Background(), "SELECT id, name, description, technologies, Start_date, End_date, image FROM tb_projects ORDER BY id DESC")
	if err != nil {
		fmt.Println("Error on SELECT : " + err.Error())
		return
	}

	var Addprj = []Prj{}

	for newprj.Next() {
		var prjData = Prj{}
		err := newprj.Scan(&prjData.ID, &prjData.PrjName, &prjData.Desc, &prjData.Tech, &prjData.Start_date, &prjData.End_date, &prjData.IMG)
		if err != nil {
			fmt.Println("Error on Scan : " + err.Error())
			return
		}

		// const format string = "2006-1-2"
		// t1, _ := time.Parse(format, prjData.Start_date)
		// t2, _ := time.Parse(format, prjData.End_date) //reflect.TypeOf(t2) = time.Time

		Duration := prjData.End_date.Sub(prjData.Start_date)
		prjData.Duration = int(Duration.Hours() / 24)

		// prjData.Format_date = prjData.Post_date.Format("2 January 2006")

		Addprj = append(Addprj, prjData)
	}

	project := map[string]interface{}{
		"Project": Addprj,
	}

	// fmt.Println(Addprj) //Result: Working

	tmpt.Execute(w, project)
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

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var Detail = Prj{}
	// disabled due to database connection
	// for index, data := range Addprj {
	// 	if index == id {
	// 		Detail = Prj{
	// 			PrjName:  data.PrjName,
	// 			Duration: data.Duration,
	// 			Desc:     data.Desc,
	// 			Tech:     data.Tech,
	// 			// IMG:   data.IMG,
	// 		}
	// 	}
	// }
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, description, technologies, Start_date, End_date, image FROM tb_projects WHERE id=$1", id).Scan(
		&Detail.ID, &Detail.PrjName, &Detail.Desc, &Detail.Tech, &Detail.Start_date, &Detail.End_date, &Detail.IMG,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	Duration := Detail.End_date.Sub(Detail.Start_date)
	Detail.Duration = int(Duration.Hours() / 24)
	Detail.Str_sdate = Detail.Start_date.Format("2 January 2006")
	Detail.Str_edate = Detail.End_date.Format("2 January 2006")

	article := map[string]interface{}{
		"Article": Detail,
	}

	tmpt.Execute(w, article)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Addprj = append(Addprj[:index], Addprj[index+1:]...)
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error on DELETE : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("register.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

func reginput(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	name := r.Form.Get("nama")
	email := r.Form.Get("surel")

	password := r.Form.Get("sandi")
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error on INSERT : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")

	session.Save(r, w)
	session.AddFlash("Registration success!", "message")

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("login.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

func loginput(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")

	email := r.Form.Get("surel")
	password := r.Form.Get("sandi")
	fmt.Println(email)

	user := User{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT name, email, password FROM tb_user WHERE email=$1", email).Scan(&user.Name, &user.Email, &user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error on SELECT : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error on HASH COMPARING : " + err.Error()))
		return
	}

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Options.MaxAge = 0

	session.AddFlash("Login success!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")
	session.Options.MaxAge = -1

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
