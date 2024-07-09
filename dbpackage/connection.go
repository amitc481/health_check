// Package dbpackage SPDX-FileCopyrightText: (C) 2022 Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel
package dbpackage

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient *mongo.Client Initialized
var MongoClient *mongo.Client

// MongoCtx context.Context Initialized
var MongoCtx context.Context

// MongoErr error Initialized
var MongoErr error

// DB name declaration
var DB string

// DbConnect - DataBase Connection
func DbConnect() {
	ctx := context.TODO()
	DbName := os.Getenv("DB_Name")
	uri := os.Getenv("DB_Write_Connection")

	opts := options.Client().ApplyURI(uri)

	// Create a new client and connect to the server
	// client, err := mongo.NewClient(opts)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
		// panic(err)
		MongoClient = nil
		MongoCtx = nil
		MongoErr = err
		// return nil, nil, err
	}
	// To access mongo client across packages
	MongoClient = client
	MongoCtx = ctx
	MongoErr = err
	DB = DbName
	// return ctx, client, nil
}

// InsertCollectionData - insert data

func GetCollectionData[T any](ctx context.Context, client *mongo.Client, databaseName string, collectionName string, filter interface{}, resultSlicePtr *[]T) error {
	// Connect to DB
	database := client.Database(databaseName)
	collection := database.Collection(collectionName)

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Println("Error in getting data from collection....")
		fmt.Println(err)
		return err
	}

	defer cursor.Close(context.TODO())

	// Decode the cursor result into resultSlicePtr
	if err = cursor.All(context.TODO(), resultSlicePtr); err != nil {
		log.Fatal(err)
	}
	return nil
}

func InsertCollectionData[T any](ctx context.Context, client *mongo.Client, databaseName string, collectionName string, insertData T) (*mongo.InsertOneResult, error) {
	//func InsertCollectionData[T any](ctx context.Context, client *mongo.Client, databaseName string, collectionName string, insertData T) (*mongo.DeleteResult, error) {
	// Connect to DB
	database := client.Database(databaseName)

	// Access specific collection
	//collection := dbpackage.MongoDatabase.Collection(table)
	collection := database.Collection(collectionName)

	// Read all data from collection
	result, err := collection.InsertOne(context.Background(), insertData)
	if err != nil {
		fmt.Println("Error in Posting data from collection....")
		fmt.Println(err)
		return result, err
	}

	return result, err
}

// DeleteCollectionData - delete data
func DeleteCollectionData(ctx context.Context, client *mongo.Client, databaseName string, collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	// Connect to DB
	database := client.Database(databaseName)

	// Access specific collection
	//collection := dbpackage.MongoDatabase.Collection(table)
	collection := database.Collection(collectionName)

	// Read all data from collection
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("Error in getting data from collection....")
		fmt.Println(err)
		return result, err
	}

	return result, err
}

// GetCollectionCount - get count
func GetCollectionCount(ctx context.Context, client *mongo.Client, databaseName string, collectionName string, filter interface{}) (int64, error) {
	database := client.Database(databaseName)
	collection := database.Collection(collectionName)
	count, err := collection.CountDocuments(ctx, filter)
	fmt.Print("In DB Lib")
	fmt.Print(count)
	if err != nil {
		fmt.Println("Error finding documents:", err)
		// Handle error
		return 0, err
	}

	return count, nil
}
