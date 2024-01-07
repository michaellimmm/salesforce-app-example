package models

import (
	"context"
	"errors"
	"github/michaellimmm/salesforce-app-example/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	AccountCollection = "account"
)

type AccountStatus string

const (
	AccountStatusUnspecified AccountStatus = "UNSPECIFIED"
	AccountStatusCreated     AccountStatus = "CREATED"
	AccountStatusLinked      AccountStatus = "LINKED"
	AccountStatusUnlinked    AccountStatus = "UNLINKED"
)

type Account struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	AccessToken       string             `bson:"access_token,omitempty"`
	RefreshToken      string             `bson:"refresh_token,omitempty"`
	InstanceUrl       string             `bson:"instance_url,omitempty"`
	ClientID          string             `bson:"client_id"`
	ClientSecret      string             `bson:"client_secret"`
	Status            string             `bson:"token_status,omitempty"`
	OrgID             string             `bson:"org_id"`
	SubscribedObjects string             `bson:"subscribed_objects,omitempty"`
	CreatedAt         time.Time          `bson:"created_at,omitempty"`
	UpdatedAt         time.Time          `bson:"updated_at,omitempty"`
	DeletedAt         *time.Time         `bson:"deleted_at,omitempty"`
}

func (a *Account) getCollection() db.CollectionProvider {
	return db.Datastore.Collection(AccountCollection)
}

func (a *Account) Save(ctx context.Context) error {
	now := time.Now()
	a.ID = primitive.NewObjectID()
	a.CreatedAt = now
	a.UpdatedAt = now

	_, err := a.getCollection().InsertOne(ctx, a)
	return err
}

func (a *Account) Update(ctx context.Context) error {
	filter := createFilter()
	filter["_id"] = a.ID

	a.UpdatedAt = time.Now()
	_, err := a.getCollection().UpdateOne(ctx, filter, bson.M{"$set": a})
	return err
}

func (a *Account) FindByID(ctx context.Context) error {
	filter := createFilter()
	filter["_id"] = a.ID

	result := a.getCollection().FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return ErrDataNotFound
		}

		return result.Err()
	}

	return result.Decode(a)
}

func (a *Account) FindByInstanceUrl(ctx context.Context) error {
	filter := createFilter()
	filter["instance_id"] = a.InstanceUrl

	result := a.getCollection().FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return ErrDataNotFound
		}

		return result.Err()
	}

	return result.Decode(a)
}

func (a *Account) FindByClientID(ctx context.Context) error {
	filter := createFilter()
	filter["$and"] = []bson.M{
		{
			"client_id": a.ClientID,
		},
	}

	result := a.getCollection().FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return ErrDataNotFound
		}

		return result.Err()
	}

	return result.Decode(a)
}

func (a *Account) FindAllByStatus(ctx context.Context, status AccountStatus) ([]Account, error) {
	filter := createFilter()
	filter["token_status"] = string(status)

	cursor, err := a.getCollection().Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []Account{}, ErrDataNotFound
		}

		return []Account{}, err
	}

	var result []Account
	if err = cursor.All(ctx, &result); err != nil {
		return []Account{}, err
	}

	return result, nil
}

func (a *Account) IsEmpty() bool {
	return a.isIDEmpty()
}

func (a *Account) isIDEmpty() bool {
	return a.ID == primitive.NilObjectID
}
