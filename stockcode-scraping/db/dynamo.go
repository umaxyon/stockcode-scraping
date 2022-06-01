package db

import (
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
	expression.Set(expression.Name("code"), expression.Value("1111"))

}
