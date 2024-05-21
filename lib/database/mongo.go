// Package database
// Description: MongoDB connection and collection initialization.
package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoIndex struct {
	Keys   interface{}
	Unique bool
}

func MongoConnect(uri, dbname string, timeout time.Duration) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()
	client, err := mongo.Connect(ctx,
		options.Client().
			SetMaxConnIdleTime(0).
			SetMinPoolSize(0).
			SetMaxPoolSize(0).
			SetMaxConnecting(100).
			ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	db := client.Database(dbname)
	if err := client.Ping(ctx, db.ReadPreference()); err != nil {
		return nil, err
	}
	return db, nil
}

func MongoInit(db *mongo.Database, collectionName string, index ...MongoIndex) *mongo.Collection {
	var (
		ctx        context.Context
		cancel     context.CancelFunc
		collection *mongo.Collection

		indexes []mongo.IndexModel
	)

	ctx, cancel = context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	var collectionValidate = func() (created bool) {
		list, err := db.ListCollectionNames(ctx, bson.M{})
		if err != nil {
			println("[MONGO-DB-] ListCollectionNames error:", err.Error())
			return
		}
		for _, n := range list {
			if n == collectionName {
				created = true
				break
			}
		}
		return
	}
	if created := collectionValidate(); created {
		println("[MONGO-DB-] Collection", collectionName, "is already available. collectionValidate=", created)
	} else {
		if err := db.CreateCollection(ctx, collectionName); err != nil {
			println("[MONGO-DB-] CreateCollection error:", err.Error())
		} else {
			println("[MONGO-DB-] Collection", collectionName, "is already available")
		}
	}
	collection = db.Collection(collectionName)

	for _, uq := range index {
		if uq.Keys == nil {
			continue
		}
		indexes = append(indexes, mongo.IndexModel{
			Keys:    uq.Keys,
			Options: options.Index().SetUnique(uq.Unique),
		})
	}

	if len(indexes) > 0 {
		names, err := collection.Indexes().CreateMany(ctx, indexes)
		if err != nil {
			println("[MONGO-DB-] Make indexes error:", err.Error())
		}
		for _, name := range names {
			println("[MONGO-DB-] Index created", name)
		}
	}

	return collection
}
