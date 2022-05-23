package yh

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"stockcode-scraping/lib"
	"strings"
	"sync"
)

type StockCodeLink struct {
	Code, Name, Topic string
}

func NewStockCodeLink() *StockCodeLink {
	return &StockCodeLink{}
}

type Ret struct {
	chunkStockCodeList []StockCodeLink
	errList            []error
}

func GetAllStockCodeLinkList(c chan<- Ret, industryChunk []IndustryLink) {
	chunkStockCodeList, errList :=
		lib.LoopScraping[IndustryLink, StockCodeLink](industryChunk, ScrapingStockCodeLinkList)
	c <- Ret{chunkStockCodeList, errList}
}

func ScrapingStockPageStart(industryLinkList []IndustryLink) error {
	industryChunkList := lib.ListChunk(industryLinkList, ThreadCount)

	var allList []StockCodeLink
	var allErrList []error
	var wg sync.WaitGroup
	ch := make(chan Ret)
	for i := range industryChunkList {
		wg.Add(1)
		go GetAllStockCodeLinkList(ch, industryChunkList[i])

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

func ScrapingStockCodeLinkList(indLink IndustryLink) ([]StockCodeLink, error) {
	doc, err := lib.GetPageDocument(indLink.Url)
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

			nextScList, ee := ScrapingStockCodeLinkList(*next)
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
