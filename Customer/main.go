package Customer

import (
	datastoreCustomer "Customer/datastore"
	"Customer/driver"
	handlerCustomer "Customer/handler"
	serviceCustomer "Customer/service"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db, err := driver.ConnectToSQL()
	if err != nil {
		log.Println("could not connect to sql, err:", err)
		return
	}

	customerStore := datastoreCustomer.New(db)

	svcCustomer := serviceCustomer.New(customerStore)
	customer := handlerCustomer.New(svcCustomer)

	r := mux.NewRouter()

	r.HandleFunc("/customer", customer.Post).Methods(http.MethodPost)

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	fmt.Println("Success: server is stared")

	err = server.ListenAndServe()

	if err != nil {
		fmt.Println(err)
	}
}
