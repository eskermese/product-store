package storage

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	Product *Product
}

func New(db *mongo.Database) *Storage {
	return &Storage{
		Product: NewProduct(db.Collection("products")),
	}
}
