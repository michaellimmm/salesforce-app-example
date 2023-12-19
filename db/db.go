package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Datastore DBProvider

type (
	DBProvider interface {
		Collection(name string) CollectionProvider
		Disconnect(ctx context.Context) error
		SelectDB(dbName string)
		DropDB(ctx context.Context) error
		GetClient() *mongo.Client
		TransactionWrapper
	}

	TransactionWrapper interface {
		WithTransaction(
			context.Context,
			func(context.Context) (interface{}, error),
			...*options.TransactionOptions,
		) (interface{}, error)
	}

	DBOptions func(c *options.ClientOptions)

	DB struct {
		client *mongo.Client
		db     *mongo.Database
	}
)

func NewDB(
	ctx context.Context,
	opts ...DBOptions,
) DBProvider {

	clientOptions := options.Client()
	for _, o := range opts {
		o(clientOptions)
	}

	if err := clientOptions.Validate(); err != nil {
		log.Fatal("invalid client options: %w", err)
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("failed to connect to mongodb server: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("failed to ping mongodb server: %w", err)
	}

	return &DB{client: client}
}

func WithURI(uri string) DBOptions {
	return func(c *options.ClientOptions) {
		c.ApplyURI(uri)
	}
}

func WithPlainAuth(user, pass string) DBOptions {
	return func(c *options.ClientOptions) {
		c.SetAuth(options.Credential{
			AuthMechanism: "PLAIN",
			Username:      user,
			Password:      pass,
		})
	}
}

func WithMaxPoolSize(size uint64) DBOptions {
	return func(c *options.ClientOptions) {
		c.SetMaxPoolSize(size)
	}
}

func WithMinPoolSize(size uint64) DBOptions {
	return func(c *options.ClientOptions) {
		c.SetMinPoolSize(size)
	}
}

func (d *DB) Collection(name string) CollectionProvider {
	return &baseCollection{d.db.Collection(name)}
}

func (d *DB) Disconnect(ctx context.Context) error {
	if d.client == nil {
		return nil
	}

	return d.client.Disconnect(ctx)
}

func (d *DB) GetClient() *mongo.Client {
	return d.client
}

func (d *DB) SelectDB(dbName string) {
	if d.client == nil {
		return
	}

	d.db = d.client.Database(dbName)
}

func (d *DB) DropDB(ctx context.Context) error {
	return d.db.Drop(ctx)
}

func (d *DB) WithTransaction(
	ctx context.Context,
	txFunc func(context.Context) (interface{}, error),
	opts ...*options.TransactionOptions,
) (interface{}, error) {
	session, err := d.client.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}

	defer session.EndSession(ctx)

	return session.WithTransaction(
		ctx,
		func(sCtx mongo.SessionContext) (interface{}, error) { return txFunc(sCtx) },
		opts...,
	)
}
