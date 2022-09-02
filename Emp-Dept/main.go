package main

import (
	storeDept "developer.zopsmart.com/go/gofr/Emp-Dept/datastore/department"
	storeEmp "developer.zopsmart.com/go/gofr/Emp-Dept/datastore/employee"
	handlerDept "developer.zopsmart.com/go/gofr/Emp-Dept/handler/department"
	handlerEmp "developer.zopsmart.com/go/gofr/Emp-Dept/handler/employee"
	serviceDept "developer.zopsmart.com/go/gofr/Emp-Dept/service/department"
	serviceEmp "developer.zopsmart.com/go/gofr/Emp-Dept/service/employee"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.Server.ValidateHeaders = false
	// enabling /swagger endpoint for Swagger UI
	app.EnableSwaggerUI()

	// employee
	empStore := storeEmp.New()
	empService := serviceEmp.New(empStore)
	empHandler := handlerEmp.New(empService)

	// department
	deptStore := storeDept.New()
	deptService := serviceDept.New(deptStore)
	deptHandler := handlerDept.New(deptService)

	app.POST("/employee", empHandler.Post)

	app.PUT("/employee/{id}", empHandler.Put)

	app.DELETE("/employee/{id}", empHandler.Delete)

	app.GET("/employee/{id}", empHandler.Get)

	app.GET("/employee", empHandler.GetAll)

	app.POST("/department", deptHandler.Post)

	app.PUT("/department/{id}", deptHandler.Put)

	app.DELETE("/department/{id}", deptHandler.Delete)

	// starts the server
	app.Start()
}
