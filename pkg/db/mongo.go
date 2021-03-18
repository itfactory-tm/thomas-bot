package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	conn *mongo.Database
}

// TODO: Add cache!
func (m *MongoDatabase) ConfigForGuild(guildID string) (*Configuration, error) {
	var conf Configuration
	err := m.conn.Collection("configuration").FindOne(context.TODO(), bson.D{{"guildID", guildID}}).Decode(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func NewMongoDB(url, db string) (Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	return &MongoDatabase{
		conn: client.Database(db),
	}, nil
}
