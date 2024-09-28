package main

import (
	"context"
	"fmt"
	"os"
	"serverless-aws-cdk/internal/db"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// ScanItems fetches all items from a DynamoDB table
func WriteItems(tableName string) ([]map[string]*dynamodb.AttributeValue, error) {
	// Add logic to scan items from DynamoDB
	db := db.NewDB()
	ctx := context.Background()

	items, err := db.ScanItems(ctx, tableName)

	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	return items, nil
}

func main() {
	tableName := os.Getenv("TABLE_NAME")

	items, err := WriteItems(tableName)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Items: ", items)
}
