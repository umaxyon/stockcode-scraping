package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"stockcode-scraping/yh"
)

func handler() error {
	var err error = nil
	ind := yh.NewIndustry()

	err = ind.GetIndustryLinkList()
	if err != nil {
		return err
	}

	err = yh.ScrapingStockPageStart(ind.LinkList)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
