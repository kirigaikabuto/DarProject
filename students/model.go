package students
type Repository interface {
	GetStudents() ([]*Student,error)
	AddStudent(st *Student) (*Student,error)
}


type Student struct{
	Id int64 `json:"id,pk"`
	FirstName string `json:"firstname,omitempty"`
	LastName string `json:"lastname,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}