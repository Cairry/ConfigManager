package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	GvaMongoCli         *mongo.Client
	GvaMongoCtx         context.Context
	ConfigCollection    *mongo.Collection
	HisConfigCollection *mongo.Collection
)

func init() {

	uri := "mongodb://root:123456@192.168.1.158:27017"

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(GvaMongoCtx, clientOptions)
	if err != nil {
		log.Println(err)
		return
	}

	err = client.Ping(GvaMongoCtx, nil)
	if err != nil {
		log.Println("Connect to MongoDB failed, err:", err)
		return
	}
	log.Println("Connect to MongoDB success")

	GvaMongoCli = client

	initDatabases()

}

func initDatabases() {

	ConfigCollection = GvaMongoCli.Database("semaik").Collection("config")
	HisConfigCollection = GvaMongoCli.Database("semaik").Collection("hisConfig")

}
