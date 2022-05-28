package db

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoAccessor struct {
	db *dynamodb.DynamoDB
}

func NewAccessor() *DynamoAccessor {
	ses, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	db := dynamodb.New(ses)
	return &DynamoAccessor{db}
}
