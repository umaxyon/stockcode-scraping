package main

import (
	"fmt"
	"stockcode-scraping/db"
	"stockcode-scraping/db/test"
	"testing"
)

func TestMain(m *testing.M) {
	test.PrepareTestAspect(func() int {
		return m.Run()
	}, "../dynamo_ddl.yaml")
}

func TestHandler(t *testing.T) {
	//t.Run("Unable to get IP", func(t *testing.T) {
	//	yh.YahooProfileUrl = "https://nothing"
	//	err := handler()
	//	if err == nil {
	//		t.Fatal("Error failed to trigger with an invalid request")
	//	}
	//})
	//
	//t.Run("Non 200 Response", func(t *testing.T) {
	//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		w.WriteHeader(500)
	//	}))
	//	defer ts.Close()
	//
	//	yh.YahooProfileUrl = ts.URL
	//
	//	err := handler()
	//	if err != nil && err.Error() != yh.ErrNon200Response.Error() {
	//		t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
	//	}
	//})

	t.Run("Successful Request", func(t *testing.T) {
		err := handler()
		if err != nil {
			t.Fatal("Everything should be ok")
		}
		ret := db.NewAccessor().Query()

		fmt.Println(ret)

	})
}
