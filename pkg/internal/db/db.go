package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joho/godotenv"
)

// DB is a struct that holds the DynamoDB client
type DB struct {
	client *dynamodb.DynamoDB
}

// NewDB creates a new DB instance
func NewDB() *DB {
	dotenvErr := godotenv.Load()
	if dotenvErr != nil {
		log.Fatal("Error loading .env file")
	}

	var sess *session.Session
	var err error

	environment := os.Getenv("ENVIRONMENT")
	if environment == "local-db" {
		sess, err = session.NewSession(&aws.Config{
			Region:   aws.String("us-east-1"),
			Endpoint: aws.String("http://127.0.0.1:8000"),
		})
	} else {
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"",
			),
		})
	}

	if err != nil {
		panic(fmt.Errorf("failed to create session: %w", err))
	}

	log.Println("Successfully created session")

	return &DB{
		client: dynamodb.New(sess),
	}
}

// GetItem fetches an item from DynamoDB
func (db *DB) GetItem(ctx context.Context, tableName, id string) (map[string]*dynamodb.AttributeValue, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	result, err := db.client.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return result.Item, nil
}

// PutItem inserts an item into DynamoDB
func (db *DB) PutItem(ctx context.Context, tableName string, item map[string]*dynamodb.AttributeValue) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	_, err := db.client.PutItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

// DeleteItem deletes an item from DynamoDB
func (db *DB) DeleteItem(ctx context.Context, tableName, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	_, err := db.client.DeleteItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

// UpdateItem updates an item in DynamoDB
func (db *DB) UpdateItem(ctx context.Context, tableName string, id string, updateExpression string, expressionAttributeValues map[string]*dynamodb.AttributeValue) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	_, err := db.client.UpdateItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

// QueryItems queries items from DynamoDB
func (db *DB) QueryItems(ctx context.Context, tableName string, keyConditionExpression string, expressionAttributeValues map[string]*dynamodb.AttributeValue) ([]map[string]*dynamodb.AttributeValue, error) {
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(keyConditionExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	result, err := db.client.QueryWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}

	return result.Items, nil
}

// ScanItems scans items from DynamoDB
func (db *DB) ScanItems(ctx context.Context, tableName string) ([]map[string]*dynamodb.AttributeValue, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := db.client.ScanWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items: %w", err)
	}

	return result.Items, nil
}

// BatchGetItems fetches multiple items from DynamoDB
func (db *DB) BatchGetItems(ctx context.Context, tableName string, keys []map[string]*dynamodb.AttributeValue) ([]map[string]*dynamodb.AttributeValue, error) {
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: {
				Keys: keys,
			},
		},
	}

	result, err := db.client.BatchGetItemWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to batch get items: %w", err)
	}

	return result.Responses[tableName], nil
}

// BatchWriteItems writes multiple items to DynamoDB
func (db *DB) BatchWriteItems(ctx context.Context, tableName string, items []map[string]*dynamodb.AttributeValue) error {
	var writeRequests []*dynamodb.WriteRequest
	for _, item := range items {
		writeRequests = append(writeRequests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: item,
			},
		})
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			tableName: writeRequests,
		},
	}

	_, err := db.client.BatchWriteItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to batch write items: %w", err)
	}

	return nil
}
