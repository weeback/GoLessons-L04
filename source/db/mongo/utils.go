package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	RoleReadOnly  dbRole = "readOnly"
	RoleReadWrite dbRole = "readWrite"
)

type dbRole string

func (val dbRole) validate() error {
	switch val {
	case RoleReadOnly, RoleReadWrite:
		return nil
	default:
		return errors.New("invalid role")
	}
}

// validationRoles validate roles
// and check if the roles have the correct permissions
// to read and write to the database
// return true if all roles are valid and have the correct permissions
// otherwise return false
func validationRoles(db *mongo.Database, roles ...dbRole) bool {
	var (
		randomCollection = "random-" + uuid.NewString()
		sampleCollection = db.Collection(randomCollection)
		validatePassed   = true
	)
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)

	defer func() {
		// try to drop sample collection
		if err := sampleCollection.Drop(ctx); err != nil {
			fmt.Printf("Unable to drop collection: %s - %s\n", randomCollection, err.Error())
		}
		// cancel context
		cancel()
	}()

	//
	for _, role := range roles {
		if err := role.validate(); err != nil {
			panic(fmt.Sprintf("Invalid role: %s", role))
		}
		// check read permission
		if role == RoleReadOnly || role == RoleReadWrite {
			r := sampleCollection.FindOne(ctx, bson.M{})
			if err := r.Err(); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				fmt.Printf("Permision denied: Role (%s) can access to db (read) - %s\n", role, err.Error())
				validatePassed = false
			}
		}
		// check write permission
		if role == RoleReadWrite {
			_, err := sampleCollection.InsertOne(ctx, bson.M{"random_value": uuid.NewString()})
			if err != nil {
				fmt.Printf("Permision denied: Role (%s) can access to db (write) - %s\n", role, err.Error())
				validatePassed = false
			}
		}
	}
	return validatePassed
}
