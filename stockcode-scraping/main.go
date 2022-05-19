package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	YahooProfileUrl   = "https://profile.yahoo.co.jp/"
	ErrNon200Response = errors.New("Non 200 Response found")
)

func handler() error {
	resp, err := http.Get(YahooProfileUrl)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return ErrNon200Response
	}

	doc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(doc))
	return nil
}

func main() {
	lambda.Start(handler)
}
