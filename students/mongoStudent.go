package students

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var (
	collection *mongo.Collection
)
type MongoConfig struct {
	Host string
	Database string
	Port string
}
type repo struct{
	dbcon *mongo.Database
}

func NewStudentRepository(config MongoConfig) (Repository,error){
	clientOptions:=options.Client().ApplyURI("mongodb://"+config.Host+":"+config.Port)
	client,err := mongo.Connect(context.TODO(),clientOptions)
	if err!=nil{
		return nil,err
	}
	err = client.Ping(context.TODO(),nil)
	if err!=nil{
		return nil,err
	}
	db:=client.Database(config.Database)
	collection=db.Collection("students")
	return &repo{dbcon:db,},nil
}

func(mpro *repo) GetStudents() ([]*Student,error){
	findOptions:=options.Find()
	var students []*Student
	cur,err :=collection.Find(context.TODO(),bson.D{{}},findOptions)
	if err!=nil{
		return nil,err
	}
	for cur.Next(context.TODO()){
		var student Student
		err:=cur.Decode(&student)
		if err!=nil{
			return nil,err
		}
		students = append(students,&student)
	}
	if err:=cur.Err();err!=nil{
		return nil,err
	}
	cur.Close(context.TODO())
	return students,nil
}
func (mpro *repo) AddStudent(st *Student) (*Student,error){
	insertResult,err:=collection.InsertOne(context.TODO(),st)
	if err!=nil{
		return nil,err
	}
	fmt.Println("Inserted document",insertResult.InsertedID)
	return st,nil
}

