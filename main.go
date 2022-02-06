package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"log"
	"mongo-golang/controllers"
	"net/http"
)

func main() {

	r := httprouter.New()

	uc := controllers.NewUserController(getSession())

	r.GET("/", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.WriteHeader(200)
		fmt.Fprint(writer, "Hello World")
	})
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user/", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)

	http.ListenAndServe("localhost:9000", r)
}

func getSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	log.Println("Connected To MongoDB")
	return s
}
