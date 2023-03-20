package mongo_client

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const DB_NAME = "ssc"
const DB_URI = "mongodb://localhost:27017"
const FILE_COLLECTION = "fileobjects"


func Connect() (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(DB_URI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return client, nil
}

func GetCollection(collectionName string, client mongo.Client) *mongo.Collection {
	return client.Database(DB_NAME).Collection(collectionName)
}

func RunQuery(id string, exts []string) ([]*SearchObject, error) {
	mongoClient, err := Connect()
	if err != nil {
		return nil, fmt.Errorf("runQuery() failed, %v", err)
	}
	objects := GetCollection(FILE_COLLECTION, *mongoClient)
	results, err := queryPath(objects, id, exts)
	if err != nil {
		return nil, fmt.Errorf("runQuery() failed to execute query, %v", err)
	}
	return results, nil
}

func RunProjectQuery(project string, offset int, limit int) ([]*SearchObject, error) {
	mongoClient, err := Connect()
	if err != nil {
		return nil, fmt.Errorf("runQuery() failed, %v", err)
	}
	objects := GetCollection(FILE_COLLECTION, *mongoClient)
	results, err := queryProject(objects, project, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("runProjectQuery() failed to execute query, %v", err)
	}
	return results, nil
}

