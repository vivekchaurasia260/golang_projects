package main

import (
	"employee-application/pkg/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	//rabbitmq.MessagingQueue()

	//rabbitmq.ReceivingQueue()

	//dummyData.ReadingFromFile()

	r := mux.NewRouter()
	routes.RegisterEmployeeRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:9011", r))

}
