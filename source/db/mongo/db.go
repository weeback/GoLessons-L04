package db

import (
	"os"
	"time"

	"myapp/lib/database"

	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type DB struct {
	UserDB
}

func New(uri, dbname string, timeout time.Duration) *DB {
	var (
		isLoading   = true
		loadBarExit = func(err error) {
			var (
				msg string
			)
			if err != nil {
				msg = "Unable to connect"
			} else {
				msg = "Connected"
			}

			isLoading = false
			time.Sleep(500 * time.Millisecond)
			print("\r[MONGO-DB-] ", msg)
			println("              \n")
			if err != nil {
				panic(err)
			}
		}
	)

	// Get uri and dbname from env
	if val := os.Getenv("MONGO_URI"); len(val) > 0 {
		uri = val
		dbname = os.Getenv("MONGO_DB")
		println("\r[MONGO-DB-] Mongo connection has set by environ variable! Using env MONGO_URI=", uri, " | MONGO_DB=", dbname)
	}
	// Get dbname from uri if not provided
	if len(dbname) == 0 {
		if uriOpt, err := connstring.ParseAndValidate(uri); err == nil {
			dbname = uriOpt.Database
		}
	}

	go func() {
		for dot := ""; isLoading; func() {
			time.Sleep(500 * time.Millisecond)
			if len(dot) < 3 {
				dot += "."
			} else {
				dot = ""
			}
		}() {
			print("\r[MONGO-DB-] [⣿⣿⣿⣿⣿⣿⣦⣿] connecting", dot, "   ")
		}
	}()

	println("[MONGO-DB-] URI=", uri, "| DbName=", dbname)
	// Connect to MongoDB
	db, err := database.MongoConnect(uri, dbname, timeout)
	if err != nil {
		println("\r[MONGO-DB-] MongoConnect to :", uri, ", dbname :", dbname, ", timeout :", timeout)
		loadBarExit(err)
		return nil
	} else {
		loadBarExit(nil)
	}
	//
	if ok := validationRoles(db, RoleReadOnly, RoleReadWrite); ok {
		println("[MONGO-DB-] Connection has been established successfully")
	} else {
		println("Warning: some collection permission is not valid")
		os.Exit(1)
	}
	// return the database instance
	return &DB{
		UserDB: NewUserDB(db),
	}
}
