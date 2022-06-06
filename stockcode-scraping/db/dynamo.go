package db

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"stockcode-scraping/db/test"
	"stockcode-scraping/lib"
	"stockcode-scraping/yh"
)

type DynamoAccessor struct {
	db *dynamodb.DynamoDB
}

func NewAccessor() *DynamoAccessor {
	ses, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	config := aws.NewConfig().WithRegion(lib.AwsRegion)
	if test.IsTest {
		config = config.WithEndpoint(test.LocalEndpoint)
	}
	db := dynamodb.New(ses, config)
	return &DynamoAccessor{db}
}

func (da *DynamoAccessor) BatchWrite(ctx context.Context, data []yh.StockPage) error {
	items := make([]*dynamodb.WriteRequest, 0, len(data))
	for _, v := range data {
		av, _ := dynamodbattribute.MarshalMap(v)
		items = append(items, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		})
	}
	for len(items) > 0 {
		out, err := da.db.BatchWriteItemWithContext(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				"stcode": items,
			},
		})
		if err != nil {
			return fmt.Errorf("batch write to %s: %w", "stockpage", err)
		}

		items = append(items[:0], out.UnprocessedItems["stcode"]...)
	}
	return nil
}

func (da *DynamoAccessor) SaveStCode(stPageList []yh.StockPage) {
	ctx := context.Background()

	for i := 0; i < len(stPageList); i += 25 {
		end := i + 25
		if end > len(stPageList) {
			end = len(stPageList)
		}
		if err := da.BatchWrite(ctx, stPageList[i:end]); err != nil {
			fmt.Printf("SaveStCode err i=%d, end=%d, %s\n", i, end, err.Error())
		}
	}

	//_, err := da.db.PutItem(&dynamodb.PutItemInput{
	//	TableName: aws.String("stcode"),
	//	Item: map[string]*dynamodb.AttributeValue{
	//		"code": {
	//			S: aws.String("1111"),
	//		},
	//		"name": {
	//			S: aws.String("あああ"),
	//		},
	//		"tel": {
	//			S: aws.String("123-4567-8901"),
	//		},
	//	},
	//})
	//if err != nil {
	//	panic(err)
	//}
}

func (da *DynamoAccessor) Query() *dynamodb.QueryOutput {
	keyCond := expression.Key("Code").Equal(expression.Value("1111"))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		panic(err)
	}
	result, err := da.db.Query(&dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 aws.String("stcode"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	return result
}
