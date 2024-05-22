package main

import (
	"CoffeeStore"
	"CoffeeStore/api"
	"CoffeeStore/source/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	bindAddr = "localhost:8082"
)

func main() {
	fmt.Printf("Server API\n------------------------------------\n"+
		"\tVersion: %s\n------------------------------------\n", CoffeeStore.Version)
	// Connect to MongoDB
	dbCtrl := db.New("mongodb://localhost:27017", "coffeeDB", 30*time.Second)
	r := gin.New()
	api.New().RegisterRouter(r, dbCtrl)

	//id, err := dbCtrl.CreateProductDB(context.TODO(), "Matcha", "From the tea green", "asdf", []string{"apple", "banana", "cherry"})
	//if err != nil {
	//	fmt.Printf("Error creating product DB: %v\n", err)
	//	return
	//}
	fmt.Println(dbCtrl)

	// Create the server configuration
	serv := http.Server{
		Addr:    bindAddr,
		Handler: r,
	}
	// Print the bind address
	fmt.Printf("Server listening on: %s ...\n", bindAddr)

	// Start the server
	if err := serv.ListenAndServe(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	log.Fatal(http.ListenAndServe(":8082", nil))
}
