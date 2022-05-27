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

func EquallyDivide(allCnt int, sep int) []int {
	div := allCnt / sep
	rem := allCnt % sep
	ret := make([]int, sep)
	for i := 0; i < sep; i++ {
		ret[i] = div
		if i < rem {
			ret[i]++
		}
	}
	return ret
}

func ListChunk[T any](linkList []T, size int) [][]T {
	sizeArr := EquallyDivide(len(linkList), size)
	var divided [][]T
	start := 0
	end := 0
	for i := 0; i < len(sizeArr); i++ {
		end = end + sizeArr[i]
		divided = append(divided, linkList[start:end])
		start = end
	}
	return divided
}

func LoopScraping[T any, R any](linkList []T, method func(T, int) ([]R, error), threadCnt int) ([]R, []error) {
	var errList []error
	var allList []R

	for i := range linkList {
		indLink := linkList[i]
		stockCodeLinkList, e := method(indLink, threadCnt)
		if e != nil {
			errList = append(errList, e)
		} else {
			//fmt.Println(stockCodeLinkList)
			allList = append(allList, stockCodeLinkList...)
		}
	}
	return allList, errList
}
