package students

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"path"
	"html/template"
)
var(
	indexTemplate = template.Must(template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "index.html")))
)
type Endpoints interface {
	GetStudents() func(w http.ResponseWriter,r *http.Request)
	AddStudent() func(w http.ResponseWriter,r *http.Request)
	Index() func(w http.ResponseWriter,r *http.Request)
}
type endpointsFactory struct {
	studentInter Repository
}
func NewEndpointsFactory(rep Repository) Endpoints{
	return &endpointsFactory{
		studentInter: rep,
	}
}
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}
func (ef *endpointsFactory) GetStudents() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		students, err := ef.studentInter.GetStudents()
		fmt.Println(students)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, "Ошибка"+err.Error())
			return
		}
		fmt.Println(students)
		newstudents, error := json.Marshal(students)
		if error != nil {
			respondJSON(w, http.StatusInternalServerError, "Error"+err.Error())
			return
		}
		respondJSON(w, http.StatusOK, string(newstudents))
	}
}

func (ef *endpointsFactory) AddStudent() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		st:=&Student{
			FirstName: "asdsada",
			LastName:  "asdsada",
			Username:  "123213",
			Password:  "adasdadsa",
		}
		st,err:=ef.studentInter.AddStudent(st)
		if err!=nil{
			log.Fatal(err)
		}
	}
}
func (ef *endpointsFactory) Index() func(w http.ResponseWriter,r *http.Request) {
	return func(w http.ResponseWriter,r *http.Request) {
		if err := indexTemplate.ExecuteTemplate(w, "layout", nil); err != nil {
			log.Println(err.Error())
		}
	}
}