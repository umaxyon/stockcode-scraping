package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"stockcode-scraping/yh"
	"time"
)

func handler() error {
	var err error = nil
	ind := yh.NewIndustry()

	now := time.Now()
	err = ind.GetIndustryLinkList()
	if err != nil {
		return err
	}

	err = yh.ScrapingStockPageStart(ind.LinkList)
	if err != nil {
		return err
	}
	fmt.Printf("%vms\n", time.Since(now).Milliseconds())
	return nil
}

func main() {
	lambda.Start(handler)
}
