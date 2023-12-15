package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type baseCollection struct {
	*mongo.Collection
}

func (b *baseCollection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) SingleResult {
	return &singleResult{b.Collection.FindOne(ctx, filter, opts...)}
}

type CollectionProvider interface {
	FindOne(ctx context.Context, filter interface{},
		opts ...*options.FindOneOptions) SingleResult
	UpdateOne(ctx context.Context, filter interface{},
		update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	ReplaceOne(ctx context.Context, filter interface{},
		replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	InsertOne(ctx context.Context, document interface{},
		opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter interface{},
		opts ...*options.FindOptions) (*mongo.Cursor, error)
	Aggregate(ctx context.Context, pipeline interface{},
		opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateByID(ctx context.Context, id interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteMany(ctx context.Context, filter interface{},
		opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}
