package dbmongo

import (
	"context"
	"fmt"
	"gopkg.in/ini.v1"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (m *Mongo) New(cfg *ini.Section) {
	var err error

	host := cfg.Key("host").MustString("localhost")
	port := cfg.Key("port").MustInt(27017)
	database := cfg.Key("database").MustString("libraries")

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
