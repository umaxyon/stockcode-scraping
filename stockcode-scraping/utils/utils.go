package utils

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
)

var (
	ErrNon200Response = errors.New("non 200 Response found")
)

func GetPageDocument(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != 200 {
		return nil, ErrNon200Response
	}

	return goquery.NewDocumentFromReader(resp.Body)
}
