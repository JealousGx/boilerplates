package db

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	var config *aws.Config

	httpClient := &http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	creds := credentials.NewStaticCredentials(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		os.Getenv("AWS_SESSION_TOKEN"),
	)

	environment := os.Getenv("ENVIRONMENT")
	if environment == "local-db" {
		config = &aws.Config{
			Region:      aws.String("us-east-1"),
			Endpoint:    aws.String("http://host.docker.internal:8000"),
			Credentials: creds,
			HTTPClient:  httpClient,
			MaxRetries:  aws.Int(1),
			// LogLevel:   aws.LogLevel(aws.LogDebug), // Uncomment for debugging
		}

	} else {
		config = &aws.Config{
			Region:      aws.String("us-east-1"),
			Credentials: creds,
			HTTPClient:  httpClient,
		}
	}

	log.Println("Successfully created DynamoDB session")

	sess := session.Must(session.NewSession())

	return &DB{
		client: dynamodb.New(sess, config),
	}
}

// GetItem fetches an item from DynamoDB
func (db *DB) GetItem(ctx context.Context, tableName string, key map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
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
func (db *DB) DeleteItem(ctx context.Context, tableName string, key map[string]*dynamodb.AttributeValue) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}

	_, err := db.client.DeleteItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

// UpdateItem updates an item in DynamoDB
func (db *DB) UpdateItem(ctx context.Context, tableName string, key map[string]*dynamodb.AttributeValue, updateExpression *string, updateExpressionNames map[string]*string, expressionAttributeValues map[string]*dynamodb.AttributeValue) error {
	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		UpdateExpression:          updateExpression,
		ExpressionAttributeNames:  updateExpressionNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}

	_, err := db.client.UpdateItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

// QueryItems queries items from DynamoDB
func (db *DB) QueryItems(ctx context.Context, tableName string, keyConditionExpression *string, expressionAttributeNames map[string]*string, expressionAttributeValues map[string]*dynamodb.AttributeValue) ([]map[string]*dynamodb.AttributeValue, error) {
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    keyConditionExpression,    // e.g. "id = :id"
		ExpressionAttributeNames:  expressionAttributeNames,  // e.g. {"#id": "id"}
		ExpressionAttributeValues: expressionAttributeValues, // e.g. {":id": {"S": "123"}}
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
