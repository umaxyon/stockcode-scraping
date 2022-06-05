package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"stockcode-scraping/db"
	"stockcode-scraping/yh"
	"strings"
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

	stockPageList, errList := yh.ScrapingStockPageStart(ind.LinkList)
	if errList != nil {
		if len(errList) != 0 {
			var errStr string
			for i := range errList {
				errStr += errList[i].Error() + ","
			}
			err = errors.New(strings.TrimRight(errStr, ","))
		}
		return err
	}

	accessor := db.NewAccessor()
	accessor.SaveStCode(stockPageList)

	fmt.Printf("%vms\n", time.Since(now).Milliseconds())
	return nil
}

func main() {
	lambda.Start(handler)
}
