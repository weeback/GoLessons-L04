package main

import (
	"context"
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
	dbCtrl := db.New("mongodb://localhost:27017", "myapp", 30*time.Second)
	id, err := dbCtrl.UserDB.CreateSampleUser(context.TODO(), "admin", "admin")
	if err != nil {
		fmt.Printf("Unable to create sample user: %s", err.Error())
		return
	}
	fmt.Printf("Sample user created with ID: %s\n", id.Hex())
}
