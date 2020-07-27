package dbmongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (m *Mongo) Find() int64 {
	var res bson.M
	var count int64

	count = 0
	model := bson.M{"key": "/books/OL17806216M"}

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Minute)
	cur, err := m.collection.Find(ctx, model)
	if err != nil {
		panic("Cannot find document: " + err.Error())
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		err := cur.Decode(&res)
		if err != nil {
			panic("Cannot decode result: " + err.Error())
		}
		count++
	}

	return count
}
