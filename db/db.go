package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	BaseRepositoryProvider interface {
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

	BaseRepository struct {
		client *mongo.Client
		db     *mongo.Database
	}
)

// NewBaseRepository connects to the datastore and returns a BaseRepository instance
func NewBaseRepository(
	ctx context.Context,
	opts ...DBOptions,
) (BaseRepositoryProvider, error) {

	clientOptions := options.Client()
	for _, o := range opts {
		o(clientOptions)
	}

	if err := clientOptions.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client options: %w", err)
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb server: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb server: %w", err)
	}

	return &BaseRepository{client: client}, nil
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

func (br *BaseRepository) Collection(name string) CollectionProvider {
	return &baseCollection{br.db.Collection(name)}
}

func (br *BaseRepository) Disconnect(ctx context.Context) error {
	if br.client == nil {
		return nil
	}

	return br.client.Disconnect(ctx)
}

func (br *BaseRepository) GetClient() *mongo.Client {
	return br.client
}

func (br *BaseRepository) SelectDB(dbName string) {
	if br.client == nil {
		return
	}

	br.db = br.client.Database(dbName)
}

func (br *BaseRepository) DropDB(ctx context.Context) error {
	return br.db.Drop(ctx)
}

func (br *BaseRepository) WithTransaction(
	ctx context.Context,
	txFunc func(context.Context) (interface{}, error),
	opts ...*options.TransactionOptions,
) (interface{}, error) {
	session, err := br.client.StartSession()
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
