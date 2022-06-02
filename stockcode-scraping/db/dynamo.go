package db

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var (
	AWS_REGION     = "ap-northeast-1"
	IS_TEST        = false
	LOCAL_ENDPOINT = "http://localhost:8000"
)

type DynamoAccessor struct {
	db *dynamodb.DynamoDB
}

func NewAccessor() *DynamoAccessor {
	ses, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	config := aws.NewConfig().WithRegion(AWS_REGION)
	if IS_TEST {
		config = config.WithEndpoint(LOCAL_ENDPOINT)
	}
	db := dynamodb.New(ses, config)
	return &DynamoAccessor{db}
}

func (da *DynamoAccessor) Upsert() {
	_, err := da.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("stcode"),
		Item: map[string]*dynamodb.AttributeValue{
			"code": {
				S: aws.String("1111"),
			},
			"name": {
				S: aws.String("あああ"),
			},
			"tel": {
				S: aws.String("123-4567-8901"),
			},
		},
	})
	if err != nil {
		panic(err)
	}
}

func (da *DynamoAccessor) Query() {
	keyCond := expression.Key("code").Equal(expression.Value("1111"))

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
}
