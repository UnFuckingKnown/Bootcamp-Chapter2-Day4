package main

import (
	"belajar-golang/connection"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"title"   : " BELAJAR BOLANG TRANS 7",
	"isLogin" : true,
}

type Project struct {
	Id             int
	Name           string
	Start_date     time.Time
	End_date       time.Time
	Description    string
	Technologies   []string
	Image          string
	NewPostdate    int
	Author         string
}

var Projects = []Project{}


func main() {

	router := mux.NewRouter()

	connection.DatabaseConnect()

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))


	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/project", project).Methods("GET")
	router.HandleFunc("/mainblog/{id}", mainblog).Methods("GET")
	router.HandleFunc("/new-blog", newblog).Methods("POST")
	router.HandleFunc("/delete/{id}", delete).Methods("GET")
	

	fmt.Println("server running on port 5000")
	http.ListenAndServe("localhost:5000", router)

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html ; Charset=utf-8")
	w.WriteHeader(http.StatusOK)
	templ, err := template.ParseFiles("html/index.html")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("message" + err.Error()))
		return
	}
	templ.Execute(w, Data)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html ; Charset=utf-8")
	w.WriteHeader(http.StatusOK)

	templ, err := template.ParseFiles("html/blog.html")
	

	rows, err := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects;" )

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var result []Project
	for rows.Next() {
		var each = Project{}

		err = rows.Scan(&each.Id,&each.Name,&each.Start_date,&each.End_date,&each.Description,&each.Technologies,&each.Image)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result = append(result, each)
	}
	
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("message  " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Data":  Data,
		"Projects": result,
	}
	
	templ.Execute(w, resp)
}

func mainblog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html ; Charset=utf-8")
	w.WriteHeader(http.StatusOK)

	templ, err := template.ParseFiles("html/mainblog.html")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("message  " + err.Error()))
		return
	}

	CatchBlog := Project{}

	for i, data := range Projects {
		if i == id {
			CatchBlog = Project{
				Name  :        data.Name,
				Start_date:    data.Start_date,
				End_date:      data.End_date,
				Description:   data.Description,
				Technologies: data.Technologies,
				Image:        data.Image,
			}
		}
	}

	var resp = map[string]interface{}{
		"Data": Data,
		"Blogs": CatchBlog,
	}

	templ.Execute(w, resp)
}

func newblog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
		return
	}

	Projectname := r.PostForm.Get("projectname")
	Description := r.PostForm.Get("description")
	StartDate := r.PostForm.Get("startDate")
	EndDate := r.PostForm.Get("endDate")

	start, _ := time.Parse("2006-01-02", StartDate)
	end, _ := time.Parse("2006-01-02", EndDate)
    diff := end.Sub(start)



	var refilData = Project{

		Name:    Projectname,
		NewPostdate: int(diff.Hours() / 24),
		Author:   "intizam",
		Description:  Description,
	}
	Projects = append(Projects, refilData)

	http.Redirect(w, r, "/project", http.StatusMovedPermanently)
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "Charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	Projects = append(Projects[:id], Projects[id+1:]...)
	http.Redirect(w, r, "/project", http.StatusMovedPermanently)
}



