package db

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/awslabs/goformation/v6"
	"github.com/awslabs/goformation/v6/cloudformation"
	"os"
	"os/exec"
	"testing"
)

var (
	WinStartCmd = []string{"/C", "docker-compose -f .\\test\\docker-compose.yml up -d"}
	WinStopCmd  = []string{"/C", "docker-compose -p test stop"}
)

func ReadDDL() *cloudformation.Template {
	tmp, err := goformation.Open("../../dynamo_ddl.yaml")
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
	config := aws.NewConfig().WithRegion("ap-northeast-1").WithEndpoint("http://localhost:8000")
	return dynamodb.New(ses, config)
}

func TestMain(m *testing.M) {
	// Before
	err := exec.Command("cmd", WinStartCmd...).Start()
	if err != nil {
		panic(err)
	}
	tmp := ReadDDL()
	db := CreateDB()
	CreateTable(tmp, db)

	code := m.Run()

	//After
	_ = exec.Command("cmd", WinStopCmd...).Start()

	os.Exit(code)
}

func TestDynamo(t *testing.T) {
	t.Run("dynamo test", func(t *testing.T) {
		fmt.Println("hoge")
	})
}
