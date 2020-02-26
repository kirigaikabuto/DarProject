package students

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
)
var(
	indexTemplate = template.Must(template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "index.html")))
)
type Endpoints interface {
	GetStudents() func(w http.ResponseWriter,r *http.Request)
	AddStudent() func(w http.ResponseWriter,r *http.Request)
	GetStudent(idParam string) func(w http.ResponseWriter,r *http.Request)
	DeleteStudent(idParam string) func(w http.ResponseWriter,r *http.Request)
	UpdateStudent(idParam string) func(w http.ResponseWriter,r *http.Request)
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

		respondJSON(w, http.StatusOK, students)
	}
}

func (ef *endpointsFactory) AddStudent() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		data,err:=ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		student:=&Student{}
		if err:= json.Unmarshal(data,&student);err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		st,err:=ef.studentInter.AddStudent(student)
		if err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,st)
	}
}
func (ef *endpointsFactory) GetStudent(idParam string) func(w http.ResponseWriter,r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars:=mux.Vars(r)
		paramid,paramerr:=vars[idParam]
		if !paramerr{
			respondJSON(w,http.StatusBadRequest,"Не был передан аргумент")
			return
		}
		id,err:=strconv.ParseInt(paramid,10,10)
		if err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		student,err:=ef.studentInter.GetStudent(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,student)
	}
}
func (ef *endpointsFactory) DeleteStudent(idParam string) func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		vars:=mux.Vars(r)
		paramid,paramerr:=vars[idParam]
		if !paramerr{
			respondJSON(w,http.StatusBadRequest,"Не был передан аргумент")
			return
		}
		id,err:=strconv.ParseInt(paramid,10,10)
		if err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		student,err:=ef.studentInter.GetStudent(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		err=ef.studentInter.DeleteStudent(student)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,"Student was deleted")
	}
}
func (ef *endpointsFactory) UpdateStudent(idParam string) func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		vars:=mux.Vars(r)
		paramid,paramerr:=vars[idParam]
		if !paramerr{
			respondJSON(w,http.StatusBadRequest,"Не был передан аргумент")
			return
		}
		id,err:=strconv.ParseInt(paramid,10,10)
		if err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		student,err:=ef.studentInter.GetStudent(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		data,err:=ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		if err:=json.Unmarshal(data,&student);err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		updated_student,err:=ef.studentInter.UpdateStudent(student)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err)
			return
		}
		respondJSON(w,http.StatusOK,updated_student)
	}
}
func (ef *endpointsFactory) Index() func(w http.ResponseWriter,r *http.Request) {
	return func(w http.ResponseWriter,r *http.Request) {
		if err := indexTemplate.ExecuteTemplate(w, "layout", nil); err != nil {
			log.Println(err.Error())
		}
	}
}