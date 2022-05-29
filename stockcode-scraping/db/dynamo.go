package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var (
	AWS_REGION = "ap-northeast-1"
)

type DynamoAccessor struct {
	db *dynamodb.DynamoDB
}

func NewAccessor(region string) *DynamoAccessor {
	ses, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	db := dynamodb.New(ses, aws.NewConfig().WithRegion(region))
	return &DynamoAccessor{db}
}

func (da *DynamoAccessor) Upsert() {
	expression.Set(expression.Name("code"), expression.Value("1111"))

}
