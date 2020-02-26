package main

import (
	"github.com/gorilla/mux"
	"ldapExample/courses"

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
	studentrepo,err:=students.NewStudentRepository(conf)
	if err!=nil{
		log.Fatal(err)
	}
	//students
	studentendpoints:=students.NewEndpointsFactory(studentrepo)
	router.Methods("GET").Path("/").HandlerFunc(studentendpoints.Index())
	router.Methods("GET").Path("/students/").HandlerFunc(studentendpoints.GetStudents())
	router.Methods("GET").Path("/students/{id}").HandlerFunc(studentendpoints.GetStudent("id"))
	router.Methods("DELETE").Path("/students/{id}").HandlerFunc(studentendpoints.DeleteStudent("id"))
	router.Methods("PUT").Path("/students/{id}").HandlerFunc(studentendpoints.UpdateStudent("id"))
	router.Methods("POST").Path("/students/").HandlerFunc(studentendpoints.AddStudent())
	//courses
	coursesrepo,err:=courses.NewCourseRepository(conf)
	if err!=nil{
		log.Fatal(err)
	}
	coursesendpoints:=courses.NewEndpointsFactory(coursesrepo)
	router.Methods("GET").Path("/courses/").HandlerFunc(coursesendpoints.GetCourses())
	router.Methods("POST").Path("/courses/").HandlerFunc(coursesendpoints.AddCourse())
	router.Methods("GET").Path("/courses/{id}").HandlerFunc(coursesendpoints.GetCourse("id"))
	router.Methods("DELETE").Path("/courses/{id}").HandlerFunc(coursesendpoints.DeleteCourse("id"))
	router.Methods("PUT").Path("/courses/{id}").HandlerFunc(coursesendpoints.UpdateCourse("id"))
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