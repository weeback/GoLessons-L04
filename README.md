# GoLessons-L04
Đặt tên sau

* Project uses some of the following sample functions and objects:

    [lib/database/mongo.go](lib/database/mongo.go)

```struct: lib/database/mongo.go
(T) MongoIndex
    (f) Keys
    (f) Unique: bool

(F) MongoConnect(uri, dbname string, timeout time.Duration) (*mongo.Database, error)
(F) MongoInit(db *mongo.Database, collectionName string, index ...MongoIndex) *mongo.Collection
```

* Init connection to MongoDB with function `db.New(uri, dbname, timeout)` from package `source/db/mongo`

    [source/db/mongo/db.go](source/db/mongo/db.go)

```struct: source/db/mongo/db.go
(T) DB
    (f) <collection-01>
    (f) <collection-02>
    (f) <collection-03>
    
(F) NewDB(uri, dbname string, timeout time.Duration) *DB
```

```text: cmd/v1/main.go
package main

import (
	"time"
	
	"myapp/source/db/mongo"
)

func main() {
  // Connect to MongoDB
  db.NewDB("mongodb://localhost:27017", "mydb", 5*time.Second)
}
```

* Build Application
```bash
go mod tidy;\
    go build -ldflags "-s -w -extldflags '-static' -X myapp.Version=beta-1.0.0" \
    -o bin/myapp \
    -trimpath cmd/v1/*.go
```

* Run built application
```bash
chmod +x ./bin/myapp;\
  ./bin/myapp
```