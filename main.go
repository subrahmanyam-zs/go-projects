package main

import (
	department3 "EmployeeDepartment/handler/department"
	employee3 "EmployeeDepartment/handler/employee"
	department2 "EmployeeDepartment/service/department"
	"EmployeeDepartment/service/employee"
	"EmployeeDepartment/store/department"
	employee2 "EmployeeDepartment/store/employee"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	mySql := MySqlConfig{"root", "localhost", "Jason@470", "3306", "go"}
	db, err := Connection(mySql)
	if err != nil {
		fmt.Println(err)
	}
	storeEmp := employee2.New(db)
	serviceEmp := employee.New(storeEmp)
	handlerEmp := employee3.New(serviceEmp)
	storeDept := department.New(db)
	serviceDept := department2.New(storeDept)
	handlerDept := department3.New(serviceDept)

	r := mux.NewRouter()
	r.HandleFunc("/employee", handlerEmp.PostHandler).Methods("POST")
	r.HandleFunc("/employee/{id}", handlerEmp.PutHandler).Methods("PUT")
	r.HandleFunc("/employee/{id}", handlerEmp.GetHandler).Methods("GET")
	r.HandleFunc("/employee/{id}", handlerEmp.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/employee", handlerEmp.GetAll).Methods("GET")
	////
	r.HandleFunc("/department", handlerDept.PostHandler).Methods("POST")
	r.HandleFunc("/department/{id}", handlerDept.PutHandler).Methods("PUT")
	r.HandleFunc("/department/{id}", handlerDept.DeleteHandler).Methods("DELETE")
	err = http.ListenAndServe(":8080", r)
	fmt.Println(err)
}
