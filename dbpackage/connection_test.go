// SPDX-FileCopyrightText: (C) 2022 Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel
package dbpackage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"bou.ke/monkey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetConnectionData() (context.Context, *mongo.Client, error) {
	ctx := context.TODO()
	uri := os.Getenv("DB_Write_Connection")
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, nil, err
	}
	return ctx, client, nil
}

func TestDbConnect(t *testing.T) {
	t.Run("ConnectionFailure", func(t *testing.T) {
		patch := monkey.Patch(mongo.Connect, func(context.Context, ...*options.ClientOptions) (*mongo.Client, error) {
			return nil, errors.New("ConnectionError")
		})
		patch2 := monkey.Patch(log.Fatal, func(...interface{}) {})
		defer patch2.Unpatch()
		defer patch.Unpatch()
		DbConnect()
	})
	t.Run("Success", func(t *testing.T) {
		DbConnect()
	})
}

func TestGetCollectionData(t *testing.T) {
	type args struct {
		ctx            context.Context
		client         *mongo.Client
		databaseName   string
		collectionName string
		filter         interface{}
		resultSlicePtr *[]interface{}
	}
	ctx, client, _ := GetConnectionData()
	databaseName := "DKAMStaging"
	collectionName := "collection_00"
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "FilterError", args: args{ctx: ctx, client: client, databaseName: databaseName, collectionName: collectionName}, wantErr: true},
		{name: "SliceError", args: args{ctx: ctx, client: client, databaseName: databaseName, collectionName: collectionName, filter: bson.D{{}}}, wantErr: true},
		{name: "Success", args: args{ctx: ctx, client: client, databaseName: databaseName, collectionName: collectionName, filter: bson.D{{}}, resultSlicePtr: &[]interface{}{}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetCollectionData(tt.args.ctx, tt.args.client, tt.args.databaseName, tt.args.collectionName, tt.args.filter, tt.args.resultSlicePtr); (err != nil) != tt.wantErr {
				t.Errorf("GetCollectionData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsertCollectionData(t *testing.T) {
	type args struct {
		ctx            context.Context
		client         *mongo.Client
		databaseName   string
		collectionName string
		insertData     interface{}
	}
	ctx, client, _ := GetConnectionData()
	databaseName := "DKAMStaging"
	collectionName := "collection_00"
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "NoInsertData", args: args{ctx: ctx, client: client, databaseName: databaseName, collectionName: collectionName}, wantErr: true},
		{name: "Success", args: args{ctx: ctx, client: client, databaseName: databaseName, collectionName: collectionName, insertData: bson.D{{Key: "id", Value: "id01"}}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := InsertCollectionData(tt.args.ctx, tt.args.client, tt.args.databaseName, tt.args.collectionName, tt.args.insertData)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertCollectionData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDeleteCollectionData(t *testing.T) {
	type args struct {
		ctx            context.Context
		client         *mongo.Client
		databaseName   string
		collectionName string
		filter         interface{}
	}
	ctx, client, _ := GetConnectionData()
	databaseName := "DKAMStaging"
	collectionName := "collection_00"
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "FilterError", args: args{ctx: ctx, client: client, databaseName: databaseName, collectionName: collectionName}, wantErr: true},
		{name: "Success", args: args{ctx: ctx, client: client, databaseName: databaseName, collectionName: collectionName, filter: bson.D{{}}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeleteCollectionData(tt.args.ctx, tt.args.client, tt.args.databaseName, tt.args.collectionName, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCollectionData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetCollectionCount(t *testing.T) {
	type args struct {
		ctx            context.Context
		client         *mongo.Client
		databaseName   string
		collectionName string
		filter         interface{}
	}
	ctx, client, _ := GetConnectionData()
	databaseName := "DKAMStaging"
	collectionName := "collection_00"
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "FilterError", args: args{ctx: ctx, client: client, databaseName: databaseName, collectionName: collectionName}, wantErr: true},
		{name: "Success", args: args{ctx: ctx, client: client, databaseName: databaseName, collectionName: collectionName, filter: bson.D{{}}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetCollectionCount(tt.args.ctx, tt.args.client, tt.args.databaseName, tt.args.collectionName, tt.args.filter)
			fmt.Println(err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCollectionCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
