package model

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrDataNotFound error = errors.New("data not found")
)

func createFilter() bson.M {
	filter := bson.M{
		"$or": []bson.M{
			{"deleted_at": bson.M{"$exists": false}},
			{"deleted_at": nil},
		},
	}
	return filter
}
