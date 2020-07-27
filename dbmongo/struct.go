package dbmongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (m *Mongo) clean() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	m.collection.Drop(ctx)

}

func (m *Mongo) New(host string) {
	var err error
	m.client, err = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", host)))
	if err != nil {
		panic("Error when opening client")
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = m.client.Connect(ctx)
	if err != nil {
		panic("Error when opening connection")
	}

	m.collection = m.client.Database("libraries").Collection("data")
	m.clean()
}

func (m *Mongo) Close() {
	m.client.Disconnect(context.TODO())
}

func (m *Mongo) Name() string {
	return "MongoDB"
}
