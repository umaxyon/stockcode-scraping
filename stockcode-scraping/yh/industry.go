package yh

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	"sync"
)

const (
	ThreadCount = 5
)

var (
	YahooProfileUrl   = "https://profile.yahoo.co.jp"
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

type IndustryLink struct {
	Name, Url string
}

func NewIndustryLink(name, url string) *IndustryLink {
	return &IndustryLink{Name: name, Url: url}
}

type Industry struct {
	linkList []IndustryLink
}

func NewIndustry() *Industry {
	return &Industry{}
}

func (m *Industry) GetIndustryLinkList() error {
	doc, err := GetPageDocument(YahooProfileUrl)
	if err != nil {
		return err
	}

	links := doc.Find("h2 + div > table a")
	links.Each(func(i int, lnk *goquery.Selection) {
		url, _ := lnk.Attr("href")
		if !strings.HasPrefix(url, YahooProfileUrl) {
			return
		}
		m.linkList = append(m.linkList, *NewIndustryLink(lnk.Text(), url))
	})

	return nil
}

func (m *Industry) LinkChunk(size int) [][]IndustryLink {
	var divided [][]IndustryLink
	for i := 0; i < len(m.linkList); i += size {
		end := i + size
		if end > len(m.linkList) {
			end = len(m.linkList)
		}
		divided = append(divided, m.linkList[i:end])
	}
	return divided
}

func (m *Industry) GetAllStockCodeLinkList() error {
	industryChunkList := m.LinkChunk(ThreadCount)
	stockCode := NewStockCode()

	type Ret struct {
		chunkStockCodeList []StockCodeLink
		errList            []error
	}
	var allList []StockCodeLink
	var allErrList []error
	var wg sync.WaitGroup
	ch := make(chan Ret)
	for i := range industryChunkList {
		wg.Add(1)
		i := i
		go func(c chan<- Ret) {
			chunkStockCodeList, errList := stockCode.Scraping(industryChunkList[i])
			c <- Ret{chunkStockCodeList, errList}
		}(ch)
	}

	go func(c <-chan Ret) {
		for r := range c {
			if r.errList != nil {
				allErrList = append(allErrList, r.errList...)
			} else {
				allList = append(allList, r.chunkStockCodeList...)
			}
			wg.Done()
		}
	}(ch)

	wg.Wait()
	close(ch)

	if allErrList != nil {
		var errStr string
		for i := range allErrList {
			errStr += allErrList[i].Error() + ","
		}
		return errors.New(strings.TrimRight(errStr, ","))
	}

	fmt.Println(fmt.Sprintf("len(allList) = %d", len(allList)))
	return nil
}

type StockCodeLink struct {
	Code, Name, Topic string
}

func NewStockCodeLink() *StockCodeLink {
	return &StockCodeLink{}
}

type StockCode struct {
	codeList []StockCodeLink
}

func NewStockCode() *StockCode {
	return &StockCode{}
}

func (sc *StockCode) Scraping(linkList []IndustryLink) ([]StockCodeLink, []error) {
	var errList []error
	var allList []StockCodeLink
	for i := range linkList {
		indLink := linkList[i]
		fmt.Println(indLink.Name)
		stockCodeLinkList, e := sc.GetStockCodeLinkList(indLink)
		if e != nil {
			errList = append(errList, e)
		} else {
			//fmt.Println(stockCodeLinkList)
			allList = append(allList, stockCodeLinkList...)
		}
	}
	return allList, errList
}

func (sc *StockCode) GetStockCodeLinkList(indLink IndustryLink) ([]StockCodeLink, error) {
	doc, err := GetPageDocument(indLink.Url)
	if err != nil {
		return nil, err
	}

	stockCodeLinkList := []StockCodeLink{*NewStockCodeLink()}
	i := 0
	doc.Find("div.profile > table:nth-of-type(2) table td").Each(func(_ int, s *goquery.Selection) {
		txt := s.Text()
		if txt != "" {
			sc := &stockCodeLinkList[len(stockCodeLinkList)-1]
			switch (i + 1) % 3 {
			case 1:
				sc.Code = txt
			case 2:
				sc.Name = txt
			case 0:
				sc.Topic = txt
				stockCodeLinkList = append(stockCodeLinkList, *NewStockCodeLink())
				//fmt.Println(i, " ", sc)
			}
			i++
		}
	})
	if stockCodeLinkList[len(stockCodeLinkList)-1].Code == "" {
		stockCodeLinkList = stockCodeLinkList[:len(stockCodeLinkList)-1]
	}

	var pagingErrList []error
	doc.Find("div.profile > table:nth-of-type(1) a").Each(func(_ int, s *goquery.Selection) {
		if s.Text() == "次の20件" {
			nextUrl, _ := s.Attr("href")
			var next *IndustryLink = &IndustryLink{}
			*next = indLink
			next.Url = fmt.Sprintf("%s%s", YahooProfileUrl, nextUrl)

			nextScList, ee := sc.GetStockCodeLinkList(*next)
			if ee != nil {
				pagingErrList = append(pagingErrList, ee)
			} else {
				stockCodeLinkList = append(stockCodeLinkList, nextScList...)
			}
		}
	})
	if len(pagingErrList) != 0 {
		var errStr string
		for i := range pagingErrList {
			errStr += pagingErrList[i].Error() + ","
		}
		err = errors.New(strings.TrimRight(errStr, ","))
	}

	return stockCodeLinkList, err
}
