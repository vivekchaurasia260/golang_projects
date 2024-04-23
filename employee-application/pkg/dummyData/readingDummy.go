package dummyData

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Employee struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Salary    string `json:"salary"`
}

func ReadingFromFile() []Employee {

	data, err := ioutil.ReadFile("C:/Workiva/development/go-code/employee-application/pkg/dummyData/jsonData.json")

	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	// Create a struct of type Employee
	//NewEmployee := []models.Employee{}
	var employees []Employee

	fmt.Println("UnMarshalling JSON from file")

	// UnMarshall the JSON data into the Employee struct
	err = json.Unmarshal(data, &employees)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	fmt.Println("Printing parsed data")

	// Print the Parse Data
	//fmt.Println(employees)

	// for _, employee := range employees {
	// 	fmt.Println(employee)
	// }

	return employees

}
