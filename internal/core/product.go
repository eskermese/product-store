package core

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ProductDefaultSort  = "name"
	ProductSortSafeList = []string{"name", "-name", "price", "-price"}
)

type Product struct {
	ID    primitive.ObjectID `bson:"id" csv:"-"`
	Name  string             `bson:"name" csv:"PRODUCT NAME"`
	Price int                `bson:"price" csv:"PRICE"`
}
