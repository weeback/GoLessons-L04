package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"myapp/lib/database"
	"time"
)

type UserDB interface {
	CreateSampleUser(ctx context.Context, username, password string) (primitive.ObjectID, error)
}

// NewUserDB create a new instance of UserDB
// and return the UserDB interface.
// Create indexes with unique values and make the search faster.
func NewUserDB(db *mongo.Database) UserDB {
	return &users{
		co: database.MongoInit(db, "users",
			// Create indexes with unique values
			database.MongoIndex{Keys: bson.D{{"username", 1}}, Unique: true},
			// Create indexes for the collection to make the search faster
			database.MongoIndex{Keys: bson.D{{"create_at", 1}}, Unique: false},
			database.MongoIndex{Keys: bson.D{{"create_at", -1}}, Unique: false},
		),
	}
}

type users struct {
	co *mongo.Collection
}

func (u *users) CreateSampleUser(ctx context.Context, username, password string) (primitive.ObjectID, error) {
	var (
		payload = map[string]interface{}{
			"username":  username,
			"password":  password,
			"create_at": time.Now(),
		}
	)
	r, err := u.co.InsertOne(ctx, payload)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return r.InsertedID.(primitive.ObjectID), nil
}
