package routes

import (
	"github.com/gorilla/mux"

	"employee-application/pkg/controllers"
)

var RegisterEmployeeRoutes = func(router *mux.Router) {
	router.HandleFunc("/employee/", controllers.PublishEmployee).Methods("POST")
	router.HandleFunc("/employee/consume/", controllers.ConsumeEmployee).Methods("POST")
	router.HandleFunc("/employee/", controllers.GetAllEmployees).Methods("GET")
	router.HandleFunc("/employee/{empId}", controllers.GetEmployeeById).Methods("Get")
	router.HandleFunc("/employee/{empId}", controllers.UpdateEmployee).Methods("PUT")
	router.HandleFunc("/employee/{empId}", controllers.DeleteEmployee).Methods("DELETE")
}
