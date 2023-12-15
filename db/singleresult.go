package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SingleResult interface {
	Decode(v interface{}) error
	DecodeBytes() (bson.Raw, error)
	Err() error
}

type singleResult struct {
	sr *mongo.SingleResult
}

func (s *singleResult) Decode(v interface{}) error {
	return s.sr.Decode(v)
}

func (s *singleResult) DecodeBytes() (bson.Raw, error) {
	return s.sr.Raw()
}

func (s *singleResult) Err() error {
	return s.sr.Err()
}
