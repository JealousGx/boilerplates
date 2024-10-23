package testTable

import (
	"context"
	"fmt"
	"serverless-aws-cdk/internal/db"
	"sync"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var database = db.NewDB()

const TABLE_NAME = "ServerlessAWSCDKLocal"

func GetItem(pk, sk string) (map[string]*dynamodb.AttributeValue, error) {
	ctx := context.TODO()

	key, err := dynamodbattribute.MarshalMap(map[string]string{
		"pk": pk,
		"sk": sk,
	})

	if err != nil {
		return nil, err
	}

	return database.GetItem(ctx, TABLE_NAME, key)
}

func PutItem(item interface{}) error {
	ctx := context.TODO()

	marshalledItem, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	return database.PutItem(ctx, TABLE_NAME, marshalledItem)
}

func DeleteItem(pk, sk string) error {
	ctx := context.TODO()

	key, err := dynamodbattribute.MarshalMap(map[string]string{
		"pk": pk,
		"sk": sk,
	})

	if err != nil {
		return err
	}

	return database.DeleteItem(ctx, TABLE_NAME, key)
}

func UpdateItem(pk, sk string, item map[string]interface{}) error {
	ctx := context.TODO()

	updateItem := expression.UpdateBuilder{}

	for k, v := range item {
		updateItem = updateItem.Set(expression.Name(k), expression.Value(v))
	}

	key, _ := dynamodbattribute.MarshalMap(map[string]string{
		"pk": pk,
		"sk": sk,
	})

	expr, _ := expression.NewBuilder().WithUpdate(updateItem).Build()

	return database.UpdateItem(ctx, TABLE_NAME, key, expr.Update(), expr.Names(), expr.Values())
}

func QueryItems(keys []map[string]interface{}) ([]map[string]*dynamodb.AttributeValue, error) {
	ctx := context.TODO()

	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided")
	}

	var expr expression.Builder
	var keyCond expression.KeyConditionBuilder

	for _, keyMap := range keys {
		if _, ok := keyMap["pk"]; !ok {
			return nil, fmt.Errorf("missing required key: pk")
		}

		keyCond = expression.Key("pk").Equal(expression.Value(keyMap["pk"]))

		for k, v := range keyMap {
			if k != "pk" {
				expr = expr.WithCondition(expression.Name(k).Equal(expression.Value(v)))
			}
		}
	}

	queryExpr, err := expr.WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	return database.QueryItems(ctx, TABLE_NAME, queryExpr.KeyCondition(), queryExpr.Names(), queryExpr.Values())
}

func ScanItems() ([]map[string]*dynamodb.AttributeValue, error) {
	ctx := context.TODO()

	return database.ScanItems(ctx, TABLE_NAME)
}

func BatchGetItems(keys []interface{}) ([]map[string]*dynamodb.AttributeValue, error) {
	ctx := context.TODO()
	const maxBatchSize = 100

	var allResults []map[string]*dynamodb.AttributeValue
	var wg sync.WaitGroup
	errCh := make(chan error, len(keys))

	for i := 0; i < len(keys); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(keys) {
			end = len(keys)
		}

		chunk := keys[i:end]
		wg.Add(1)

		go func(chunk []interface{}) {
			defer wg.Done()

			marshalledKeys := make([]map[string]*dynamodb.AttributeValue, len(chunk))
			for j, key := range chunk {
				marshalledKey, err := dynamodbattribute.MarshalMap(key)
				if err != nil {
					errCh <- err
					return
				}
				marshalledKeys[j] = marshalledKey
			}

			resp, err := database.BatchGetItems(ctx, TABLE_NAME, marshalledKeys)
			if err != nil {
				errCh <- err
				return
			}

			allResults = append(allResults, resp...) // equivalent to [...allResults, ...resp] in JavaScript
		}(chunk)
	}

	wg.Wait()
	close(errCh)

	var finalErr error
	for err := range errCh {
		if finalErr == nil {
			finalErr = err
		} else {
			finalErr = fmt.Errorf("%w; %v", finalErr, err)
		}
	}

	return allResults, finalErr
}

func BatchWriteItems(requestItems []interface{}) error {
	ctx := context.TODO()

	const MAX_BATCH_WRITE = 25

	var chunks [][]interface{}

	for i := 0; i < len(requestItems); i += MAX_BATCH_WRITE {
		end := i + MAX_BATCH_WRITE

		if end > len(requestItems) {
			end = len(requestItems)
		}

		chunks = append(chunks, requestItems[i:end])
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(chunks))

	for _, chunk := range chunks {
		wg.Add(1)

		go func(items []interface{}) {
			defer wg.Done()

			marshalledItems := make([]map[string]*dynamodb.AttributeValue, len(items))
			for i, item := range items {
				marshalledItem, err := dynamodbattribute.MarshalMap(item)
				if err != nil {
					errCh <- err
					return
				}

				marshalledItems[i] = marshalledItem
			}

			if err := database.BatchWriteItems(ctx, TABLE_NAME, marshalledItems); err != nil {
				errCh <- err
			}
		}(chunk)
	}

	wg.Wait()
	close(errCh)

	var finalErr error
	for err := range errCh {
		if finalErr == nil {
			finalErr = err
		} else {
			finalErr = fmt.Errorf("%w; %v", finalErr, err)
		}
	}

	return finalErr
}
