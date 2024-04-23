package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"

	"employee-application/pkg/dummyData"
	"employee-application/pkg/models"
	"employee-application/pkg/rabbitmq"
	"employee-application/pkg/utils"
)

var NewEmployee models.Employee

// GET ALL EMPLOYEES
func GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	newEmployee := models.GetAllEmployees()
	res, _ := json.Marshal(newEmployee)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GET EMPLOYEE BY ID
func GetEmployeeById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empId := vars["empId"]
	ID, err := strconv.ParseInt(empId, 0, 0)

	if err != nil {
		fmt.Print("error while parsing")
	}
	employeeDetails, _ := models.GetEmployeeById(ID)
	res, _ := json.Marshal(employeeDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// CREATE EMPLOYEE
func PublishEmployee(w http.ResponseWriter, r *http.Request) {
	const numWorkers = 5
	const maxJobPerWorker = 2

	// INITIALIZING A BUFFERED CHANNEL TO CONTROL NUMBER OF CONCURRENT JOBS
	jobCh := make(chan dummyData.Employee, numWorkers*maxJobPerWorker)
	// CHANNEL TO SIGNAL WHEN ALL JOBS ARE DONE
	doneCh := make(chan struct{})

	var employees []dummyData.Employee

	// Wait group to synchronize goroutines
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for employee := range jobCh {
				err := rabbitmq.MessagingQueue(employee)

				if err != nil {
					fmt.Printf("Error in publishing employee: %v\n", err)
				}
			}
		}()
	}
	// Launch goroutine to read JSON data from file
	go func() {
		defer close(jobCh)
		employees = dummyData.ReadingFromFile()

		for _, employee := range employees {
			jobCh <- employee
		}
	}()

	// Launch goroutine to publish JSON data to RabbitMQ
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	rabbitmq.MessagingQueue(employees)
	// }()

	// Wait for all worker goroutines to finish
	//wg.Wait()
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	<-doneCh

	// RETRIEVING JSON DATA FROM FILE

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Message published to RabbitMq server!!")
}

// CONSUMING MESSAGES FROM RABBITMQ
func ConsumeEmployee(w http.ResponseWriter, r *http.Request) {
	// msgs := make(chan []byte)
	// rabbitmq.ReceivingQueue(msgs)

	// go func() {
	// 	for d := range msgs {
	// 		log.Printf("Received published messages: %s", d)
	// 	}
	// }()

	// for msg := range msgs {
	// 	// Decode Json data from message body
	// 	CreateEmployee := &models.Employee{}
	// 	err := json.Unmarshal(msg, &CreateEmployee)
	// 	fmt.Println("Employee : ", CreateEmployee)
	// 	if err != nil {
	// 		fmt.Println("error while parsing")
	// 	}

	// 	//utils.ParseBody(employee, CreateEmployee)
	// 	newEmployee := CreateEmployee.CreateEmployee()

	// 	res, _ := json.Marshal(newEmployee)
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusCreated)
	// 	w.Write(res)
	// }

	// Create a channel to receive messages
	msgs := make(chan []byte)

	// Start consuming messages asynchronously
	go consumeMessages(msgs)

	// Process incoming messages
	for msg := range msgs {
		handleMessage(w, msg)
	}

	// for {
	// 	select {
	// 	case msg, ok := <-msgs:
	// 		if !ok {
	// 			// Channel closed, all messages consumed
	// 			close(msgs)
	// 			return
	// 		}

	// 		handleMessage(w, msg)
	// 	}
	// 	fmt.Println("Ending of Consumer message!!")
	// }
}

// consumeMessages starts consuming messages from RabbitMQ
func consumeMessages(msgs chan []byte) {
	rabbitmq.ReceivingQueue(msgs)

	close(msgs)
}

// handleMessage decodes and processes the incoming message
func handleMessage(w http.ResponseWriter, msg []byte) {
	// Decode JSON data from message body
	CreateEmployee := &models.Employee{}
	err := json.Unmarshal(msg, &CreateEmployee)
	if err != nil {
		log.Println("Error while parsing:", err)
		return
	}

	// Create employee
	CreateEmployee.CreateEmployee()

	if err != nil {
		log.Println("Error while marshalling response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//fmt.Println("I am above Response Writer")
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(msg)

	//fmt.Println("End of Response Writer!!")
}

// DELETE EMPLOYEE
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	empId := params["empId"]
	ID, err := strconv.ParseInt(empId, 0, 0)

	if err != nil {
		fmt.Println("error while parsing")
	}

	employeeDetails := models.DeleteEmployee(ID)
	res, _ := json.Marshal(employeeDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// UPDATE EMPLOYEE
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	var updateEmployee = &models.Employee{}
	utils.ParseBody(r, updateEmployee)
	vars := mux.Vars(r)
	empId := vars["empId"]
	ID, err := strconv.ParseInt(empId, 0, 0)

	if err != nil {
		fmt.Print("error while parsing")
	}

	employeeDetails, db := models.GetEmployeeById(ID)

	if updateEmployee.First_Name != "" {
		employeeDetails.First_Name = updateEmployee.First_Name
	}
	if updateEmployee.Last_Name != "" {
		employeeDetails.Last_Name = updateEmployee.Last_Name
	}
	if updateEmployee.Email != "" {
		employeeDetails.Email = updateEmployee.Email
	}
	if updateEmployee.Salary != "" {
		employeeDetails.Salary = updateEmployee.Salary
	}

	db.Save(&employeeDetails)

	res, _ := json.Marshal(employeeDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(res)
}
