package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/awslabs/goformation/v6"
	"github.com/awslabs/goformation/v6/cloudformation"
	"os"
	"os/exec"
	"stockcode-scraping/lib"
)

var (
	IsTest        = false
	LocalEndpoint = "http://localhost:8000"
	WinStartCmd   = []string{"/C", "docker-compose -f .\\test\\docker-compose.yml up -d"}
	WinStopCmd    = []string{"/C", "docker-compose -p test stop"}
)

func PrepareTestAspect(testFunc func() int, ddlYaml string) {
	// Before
	err := exec.Command("cmd", WinStartCmd...).Start()
	if err != nil {
		panic(err)
	}
	tmp := ReadDDL(ddlYaml)
	CreateTable(tmp, CreateDB())

	IsTest = true
	code := testFunc()

	//After
	_ = exec.Command("cmd", WinStopCmd...).Start()

	IsTest = false
	os.Exit(code)
}

func ReadDDL(ddlFilePath string) *cloudformation.Template {
	tmp, err := goformation.Open(ddlFilePath)
	if err != nil {
		panic(err)
	}
	return tmp
}

func CreateTable(tmp *cloudformation.Template, db *dynamodb.DynamoDB) {
	tbl, err := tmp.GetDynamoDBTableWithName("StCode")
	if err != nil {
		panic(err)
	}

	_, err = db.DeleteTable(&dynamodb.DeleteTableInput{TableName: tbl.TableName})

	defs := *tbl.AttributeDefinitions
	AttributeDefinitions := make([]*dynamodb.AttributeDefinition, 0)
	AttributeDefinitions = append(AttributeDefinitions, &dynamodb.AttributeDefinition{
		AttributeName: &defs[0].AttributeName,
		AttributeType: &defs[0].AttributeType,
	})

	keys := tbl.KeySchema
	KeySchema := make([]*dynamodb.KeySchemaElement, 0)
	KeySchema = append(KeySchema, &dynamodb.KeySchemaElement{
		AttributeName: &keys[0].AttributeName,
		KeyType:       &keys[0].KeyType,
	})

	ProvisionedThroughput := dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(1),
		WriteCapacityUnits: aws.Int64(1),
	}

	_, err = db.CreateTable(&dynamodb.CreateTableInput{
		TableName:             tbl.TableName,
		AttributeDefinitions:  AttributeDefinitions,
		KeySchema:             KeySchema,
		ProvisionedThroughput: &ProvisionedThroughput,
	})
	if err != nil {
		panic(err)
	}
}

func CreateDB() *dynamodb.DynamoDB {
	ses, _ := session.NewSession()
	config := aws.NewConfig().WithRegion(lib.AwsRegion).WithEndpoint(LocalEndpoint)
	return dynamodb.New(ses, config)
}
