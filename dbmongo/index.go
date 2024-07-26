package dbmongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *Mongo) IndexJSON() {
	mod := mongo.IndexModel{
		Keys: bson.M{
			"key": 1, // index in ascending order
		}, Options: nil,
	}

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Minute)
	_, err := m.collection.Indexes().CreateOne(ctx, mod)

	if err != nil {
		panic("Cannot create index: " + err.Error())
	}
}
