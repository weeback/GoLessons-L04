package db

import (
	"CoffeeStore/Models"
	"CoffeeStore/lib/database"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var collectionName = "hotCoffees"

type ProductDB interface {
	CreateProductDB(ctx context.Context, title, description, image string, ingredients []string) (primitive.ObjectID, error)
	GetProductById(id primitive.ObjectID) ([]ProductDB, error)
	GetProductByName(ctx context.Context, name string) ([]Models.Product, error)
}

func NewProductDB(db *mongo.Database) ProductDB {
	return &products{
		co: database.MongoInit(db, collectionName),
	}
}

type products struct {
	co *mongo.Collection
}

func (p *products) CreateProductDB(ctx context.Context, title, description, image string, ingredients []string) (primitive.ObjectID, error) {
	//TODO implement me
	var (
		payload = map[string]interface{}{
			"title":       title,
			"description": description,
			"image":       image,
			"ingredients": ingredients,
			"id":          GetTimestampLastFourMicro(),
		}
	)
	r, err := p.co.InsertOne(ctx, payload)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return r.InsertedID.(primitive.ObjectID), nil
}

func GetTimestampLastFourMicro() int {
	// Get the current time
	now := time.Now()

	// Extract the microseconds part
	microseconds := now.Nanosecond() / 1000

	// Get the last four digits of the microseconds part
	lastFour := microseconds % 10000

	return lastFour
}

func (p *products) GetProductById(id primitive.ObjectID) ([]ProductDB, error) {
	cursor, err := p.co.Find(context.TODO(), bson.M{"id": id})
	if err != nil {
		return nil, err
	}
	var products []ProductDB
	if err = cursor.All(context.TODO(), &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (p *products) GetProductByName(ctx context.Context, name string) ([]Models.Product, error) {

	cursor, err := p.co.Find(ctx, bson.M{"title": name})
	if err != nil {
		return nil, err
	}
	var products []Models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}
