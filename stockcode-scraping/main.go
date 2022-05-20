package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"stockcode-scraping/yh"
)

func handler() (interface{}, error) {
	list, err := yh.NewIndustry().GetIndustryDataList()
	if err != nil {
		return nil, err
	}

	fmt.Println(list)

	return list, nil
}

func main() {
	lambda.Start(handler)
}
