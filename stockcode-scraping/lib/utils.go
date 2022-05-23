package lib

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

func ListChunk[T any](linkList []T, size int) [][]T {
	var divided [][]T
	for i := 0; i < len(linkList); i += size {
		end := i + size
		if end > len(linkList) {
			end = len(linkList)
		}
		divided = append(divided, linkList[i:end])
	}
	return divided
}

func LoopScraping[T any, R any](linkList []T, method func(T) ([]R, error)) ([]R, []error) {
	var errList []error
	var allList []R

	for i := range linkList {
		indLink := linkList[i]
		stockCodeLinkList, e := method(indLink)
		if e != nil {
			errList = append(errList, e)
		} else {
			//fmt.Println(stockCodeLinkList)
			allList = append(allList, stockCodeLinkList...)
		}
	}
	return allList, errList
}
