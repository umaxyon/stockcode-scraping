package yh

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
)

var (
	YahooProfileUrl   = "https://profile.yahoo.co.jp/"
	ErrNon200Response = errors.New("non 200 Response found")
)

type IndustryData struct {
	Name, Url string
}

func NewIndustryData(name, url string) *IndustryData {
	return &IndustryData{Name: name, Url: url}
}

type Industry struct {
	dataList []IndustryData
}

func NewIndustry() *Industry {
	return &Industry{}
}

func (m *Industry) GetIndustryPageDocument() (*goquery.Document, error) {
	resp, err := http.Get(YahooProfileUrl)
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

func (m *Industry) GetIndustryDataList() ([]IndustryData, error) {
	doc, err := m.GetIndustryPageDocument()

	if err != nil {
		return nil, err
	}
	var indList []IndustryData
	links := doc.Find("h2 + div > table a")
	links.Each(func(i int, lnk *goquery.Selection) {
		url, _ := lnk.Attr("href")
		if !strings.HasPrefix(url, YahooProfileUrl) {
			return
		}
		indList = append(indList, *NewIndustryData(lnk.Text(), url))
	})

	return indList, nil
}
