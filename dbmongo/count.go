package dbmongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (m *Mongo) Count() int64 {
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Minute)
	counted, err := m.collection.CountDocuments(ctx, bson.D{}, nil)
	if err != nil {
		panic("Cannot count documents in collection: " + err.Error())
	}

	return counted
}
