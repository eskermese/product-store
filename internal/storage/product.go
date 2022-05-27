package storage

import (
	"context"

	"github.com/ernur-eskermes/product-store/internal/core"
	"github.com/ernur-eskermes/product-store/pkg/filters"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	db *mongo.Collection
}

func NewProduct(db *mongo.Collection) *Product {
	return &Product{
		db: db,
	}
}

func (r *Product) GetAll(ctx context.Context, f *filters.Filters) ([]core.Product, error) {
	opts := options.FindOptions{}
	opts.SetSkip(f.Offset())
	opts.SetLimit(f.Limit())
	opts.SetSort(bson.D{{Key: f.SortColumn(), Value: f.SortDirection()}})

	cur, err := r.db.Find(ctx, bson.M{}, &opts)
	if err != nil {
		return nil, err
	}

	products := make([]core.Product, 0)
	if err = cur.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Product) GetTotalRecords(ctx context.Context) (int64, error) {
	return r.db.CountDocuments(ctx, bson.D{})
}

func (r *Product) UpdateOrCreate(ctx context.Context, products []core.Product) error {
	models := make([]mongo.WriteModel, 0, len(products))

	for _, product := range products {
		models = append(models, mongo.NewUpdateOneModel().SetFilter(
			bson.D{{Key: "name", Value: product.Name}},
		).SetUpdate(
			bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: product.Name}, {Key: "price", Value: product.Price}}}},
		).SetUpsert(true))
	}

	_, err := r.db.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))

	return err
}
