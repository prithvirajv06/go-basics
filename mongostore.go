package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const DBName = "go_buddies_db"

var MonGoClient *mongo.Client

var DbContext *context.Context

func close(cancel context.CancelFunc) {
	defer cancel()
	defer func() {
		if err := MonGoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	fmt.Print("***Connection Closed***")
}

func connect(uri string) (context.CancelFunc, error) {
	log.Println("----Trying to conencto to DB----")
	log.Printf("%v", uri)
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)
	DbContext = &ctx
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	MonGoClient = client
	return cancel, err
}

func ping() error {
	if err := MonGoClient.Ping(*DbContext, readpref.Primary()); err != nil {
		return err
	}
	log.Println("----connected successfully----")
	return nil
}

func insertOne(collectionName string, doc interface{}) (interface{}, error) {
	collection := *MonGoClient.Database(DBName).Collection(collectionName)
	result, err := collection.InsertOne(*DbContext, doc)
	if result == nil {
		log.Fatal("Something not roight while inserting")
	}
	return result, err
}

func findOne(collectionName string, doc interface{}) interface{} {
	collection := *MonGoClient.Database(DBName).Collection(collectionName)
	result := collection.FindOne(*DbContext, doc)
	if result == nil {
		log.Fatal("Something not roight while fetching document")
	}
	err := result.Decode(doc)
	if err != nil {
		log.Fatal("Something not roight while binding to struct")
	}
	return result
}
