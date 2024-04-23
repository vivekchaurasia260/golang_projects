package models

import (
	"github.com/jinzhu/gorm"

	"employee-application/pkg/config"
)

var db *gorm.DB

type Employee struct {
	gorm.Model
	First_Name string `gorm:"" json:"firstName"`
	Last_Name  string `json:"lastName"`
	Email      string `json:"email"`
	Salary     string `json:"salary"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Employee{})
}

func (emp *Employee) CreateEmployee() {
	db.NewRecord(emp)
	db.Create(&emp)

	//return emp
}

func GetAllEmployees() []Employee {
	var Employees []Employee
	db.Find(&Employees)

	return Employees
}

func GetEmployeeById(Id int64) (*Employee, *gorm.DB) {
	var getEmployee Employee
	db := db.Where("ID=?", Id).Find(&getEmployee)

	return &getEmployee, db
}

func DeleteEmployee(Id int64) *Employee {
	var employee Employee
	db.Where("ID=?", Id).Delete(&employee)

	return &employee
}

// func (emp *Employee) UpdateEmployee(Id int64) *Employee {
// 	//
// }
