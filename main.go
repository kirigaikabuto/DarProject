package main

import (
	"github.com/gorilla/mux"

	"ldapExample/students"
	"log"
	"net/http"
	"fmt"

	"strings"
)

func main(){
	fs := noDirListing(http.FileServer(http.Dir("./public/static")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	router:=mux.NewRouter()
	conf:=students.MongoConfig{
		Host:     "localhost",
		Database: "example",
		Port:     "27017",
	}
	repo,err:=students.NewStudentRepository(conf)
	if err!=nil{
		log.Fatal(err)
	}
	endpoints:=students.NewEndpointsFactory(repo)
	router.Methods("GET").Path("/").HandlerFunc(endpoints.Index())
	router.Methods("GET").Path("/students/").HandlerFunc(endpoints.GetStudents())
	router.Methods("POST").Path("/students/").HandlerFunc(endpoints.AddStudent())
	fmt.Println("Server is running")
	http.ListenAndServe(":8080",router)
}
func noDirListing(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") || r.URL.Path == "" {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}