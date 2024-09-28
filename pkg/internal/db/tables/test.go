package testTable

import (
	"context"
	"serverless-aws-cdk/internal/db"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var database = db.NewDB()

const TABLE_NAME = "test"

func GetItem(id string) (map[string]*dynamodb.AttributeValue, error) {
	ctx := context.TODO()

	return database.GetItem(ctx, TABLE_NAME, id)
}

func PutItem(item map[string]*dynamodb.AttributeValue) error {
	ctx := context.TODO()

	return database.PutItem(ctx, TABLE_NAME, item)
}

func DeleteItem(id string) error {
	ctx := context.TODO()

	return database.DeleteItem(ctx, TABLE_NAME, id)
}

func UpdateItem(item map[string]*dynamodb.AttributeValue) error {
	ctx := context.TODO()

	return database.PutItem(ctx, TABLE_NAME, item)
}

func QueryItems(keyConditionExpression string, expressionAttributeValues map[string]*dynamodb.AttributeValue) ([]map[string]*dynamodb.AttributeValue, error) {
	ctx := context.TODO()

	return database.QueryItems(ctx, TABLE_NAME, keyConditionExpression, expressionAttributeValues)
}

func ScanItems() ([]map[string]*dynamodb.AttributeValue, error) {
	ctx := context.TODO()

	return database.ScanItems(ctx, TABLE_NAME)
}

func BatchGetItems(keys []map[string]*dynamodb.AttributeValue) ([]map[string]*dynamodb.AttributeValue, error) {
	ctx := context.TODO()

	return database.BatchGetItems(ctx, TABLE_NAME, keys)
}

func BatchWriteItems(requestItems []map[string]*dynamodb.AttributeValue) error {
	ctx := context.TODO()

	return database.BatchWriteItems(ctx, TABLE_NAME, requestItems)
}
