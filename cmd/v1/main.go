package main

import (
	"fmt"
	"myapp"
	db "myapp/source/db/mongo"
	"time"
)

var (
	bindAddr = "0.0.0.0:8080"
)

func main() {
	fmt.Printf("Server API\n------------------------------------\n"+
		"\tVersion: %s\n------------------------------------\n", myapp.Version)

	// Connect to MongoDB
	db.New("mongodb://localhost:27017", "myapp", 30*time.Second)

}
