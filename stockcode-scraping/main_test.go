package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"stockcode-scraping/yh"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("Unable to get IP", func(t *testing.T) {
		yh.YahooProfileUrl = "https://nothing"
		_, err := handler()
		if err == nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}
	})

	t.Run("Non 200 Response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		defer ts.Close()

		yh.YahooProfileUrl = ts.URL

		_, err := handler()
		if err != nil && err.Error() != yh.ErrNon200Response.Error() {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
		}
	})

	t.Run("Successful Request", func(t *testing.T) {
		list, err := handler()
		fmt.Println(list)
		if err != nil {
			t.Fatal("Everything should be ok")
		}
	})
}
