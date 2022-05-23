package yh

import (
	"github.com/PuerkitoBio/goquery"
	"stockcode-scraping/lib"
	"strings"
)

const (
	ThreadCount = 5
)

var (
	YahooProfileUrl = "https://profile.yahoo.co.jp"
)

type IndustryLink struct {
	Name, Url string
}

func NewIndustryLink(name, url string) *IndustryLink {
	return &IndustryLink{Name: name, Url: url}
}

type Industry struct {
	LinkList []IndustryLink
}

func NewIndustry() *Industry {
	return &Industry{}
}

func (m *Industry) GetIndustryLinkList() error {
	doc, err := lib.GetPageDocument(YahooProfileUrl)
	if err != nil {
		return err
	}

	links := doc.Find("h2 + div > table a")
	links.Each(func(i int, lnk *goquery.Selection) {
		url, _ := lnk.Attr("href")
		if !strings.HasPrefix(url, YahooProfileUrl) {
			return
		}
		m.LinkList = append(m.LinkList, *NewIndustryLink(lnk.Text(), url))
	})

	return nil
}
