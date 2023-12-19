package model

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
	tokenCollection = "token"
)

type TokenStatus string

const (
	TokenStatusUnspecified TokenStatus = "UNSPECIFIED"
	TokenStatusPending     TokenStatus = "PENDING"
	TokenStatusLinked      TokenStatus = "LINKED"
	TokenStatusUnlinked    TokenStatus = "UNLINKED"
)

type Token struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	AccessToken  string             `bson:"access_token,omitempty"`
	RefreshToken string             `bson:"refresh_token,omitempty"`
	InstanceUrl  string             `bson:"instance_url,omitempty"`
	ClientID     string             `bson:"client_id"`
	ClientSecret string             `bson:"client_secret"`
	TokenStatus  string             `bson:"token_status,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty"`
}

func (t *Token) Save(ctx context.Context) error {
	filter := createFilter()
	filter["$and"] = []bson.M{
		{
			"client_id":     t.ClientID,
			"client_secret": t.ClientSecret,
		},
	}

	res := t.getCollection().FindOne(ctx, filter)
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		t.ID = primitive.NewObjectID()
		return t.save(ctx)
	}

	if res.Err() != nil {
		return res.Err()
	}

	return t.update(ctx, filter)
}

func (t *Token) save(ctx context.Context) error {
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	t.TokenStatus = string(TokenStatusPending)

	_, err := t.getCollection().InsertOne(ctx, t)
	return err
}

func (t *Token) update(ctx context.Context, filter interface{}) error {
	t.UpdatedAt = time.Now()
	_, err := t.getCollection().UpdateOne(ctx, filter, bson.M{"$set": t})
	return err
}

func (t *Token) FindByID(ctx context.Context) error {
	filter := createFilter()
	filter["_id"] = t.ID

	result := t.getCollection().FindOne(ctx, filter)
	if result.Err() != nil {
		return result.Err()
	}

	return result.Decode(t)
}

func (t *Token) FindByInstanceUrl(ctx context.Context) error {
	filter := createFilter()
	filter["instance_id"] = t.InstanceUrl

	result := t.getCollection().FindOne(ctx, filter)
	if result.Err() != nil {
		return result.Err()
	}

	return result.Decode(t)
}

func (t *Token) FindAllByStatus(ctx context.Context, status TokenStatus) ([]Token, error) {
	filter := createFilter()
	filter["token_status"] = string(status)

	cursor, err := t.getCollection().Find(ctx, filter)
	if err != nil {
		return []Token{}, err
	}

	var result []Token
	if err = cursor.All(ctx, &result); err != nil {
		return []Token{}, err
	}

	return result, nil
}

func (t *Token) isIDEmpty() bool {
	return t.ID == primitive.NilObjectID
}

func (t *Token) IsEmpty() bool {
	return t.isIDEmpty()
}

func (t *Token) getCollection() db.CollectionProvider {
	return db.Datastore.Collection(tokenCollection)
}

func createFilter() bson.M {
	filter := bson.M{
		"$or": []bson.M{
			{"deleted_at": bson.M{"$exists": false}},
			{"deleted_at": nil},
		},
	}
	return filter
}
