package dbmongo

import (
	"context"
	"dbtest/cli"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mongo struct {
	client     *mongo.Client
	collection *mongo.Collection
	url        string
}

func (m *Mongo) clean() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	m.collection.Drop(ctx)

}

func (m *Mongo) New(cli *cli.CLI) {
	var err error

	host := cli.MongoDB.Host
	port := cli.MongoDB.Port
	database := cli.MongoDB.Database

	m.url = fmt.Sprintf("mongodb://%s:%d", host, port)

	m.client, err = mongo.NewClient(options.Client().ApplyURI(m.url))
	if err != nil {
		panic("Error when opening client")
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = m.client.Connect(ctx)
	if err != nil {
		panic("Error when opening connection")
	}

	m.collection = m.client.Database(database).Collection("data")
	m.clean()
}

func (m *Mongo) Close() {
	m.client.Disconnect(context.TODO())
}

func (m *Mongo) Name() string {
	return "MongoDB"
}

func (m *Mongo) Url() string {
	return m.url
}
