package yh

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"stockcode-scraping/lib"
	"strings"
	"sync"
)

var (
	StockPageUrl = fmt.Sprintf("%s%s", YahooProfileUrl, "/fundamental/?s=")
)

type StockCodeLink struct {
	Code, Name, Topic string
}

func NewStockCodeLink() *StockCodeLink {
	return &StockCodeLink{}
}

type StockCodeLinkContainer struct {
	chunkStockCodeList []StockCodeLink
	errList            []error
}

type StockPage struct {
	Code, Name, Topic, Address, Tel, Market, Unit string
}

func NewStockPage() *StockPage {
	return &StockPage{}
}

type StockPageContainer struct {
	stockPageList []StockPage
	errList       []error
}

func GetAllStockCodeLinkList(c chan<- StockCodeLinkContainer, industryChunk [][]IndustryLink) {
	defer close(c)
	for i := range industryChunk {
		chunkStockCodeList, errList :=
			lib.LoopScraping[IndustryLink, StockCodeLink](industryChunk[i], ScrapingStockCodeLinkList, i)
		c <- StockCodeLinkContainer{chunkStockCodeList, errList}
	}
}

func ParallelPageScraping(
	linkList []StockCodeLink,
	method func(StockCodeLink, int) ([]StockPage, error),
	threadCnt int) ([]StockPage, []error) {

	var errList []error
	var allList []StockPage
	chunks := lib.ListChunk(linkList, ThreadCount)
	type S = struct {
		ret     []StockPage
		errList []error
	}
	var wg sync.WaitGroup
	ch := make(chan S, len(chunks))

	for i := range chunks {
		wg.Add(1)
		go func(inList []StockCodeLink) {
			defer wg.Done()
			ret, errList := lib.LoopScraping(inList, method, threadCnt)
			ch <- S{ret, errList}
		}(chunks[i])
	}
	wg.Wait()
	close(ch)

	for r := range ch {
		if r.errList != nil {
			errList = append(errList, r.errList...)
		} else {
			allList = append(allList, r.ret...)
		}
	}
	return allList, errList
}

func GetAllStockPage(c1 <-chan StockCodeLinkContainer, c2 chan<- StockPageContainer) {
	defer close(c2)
	i := 0
	for stockCodeLinkContainer := range c1 {
		chunkStockPage, errList :=
			ParallelPageScraping(stockCodeLinkContainer.chunkStockCodeList, ScrapingStockPage, i)
		c2 <- StockPageContainer{chunkStockPage, errList}
		i++
	}
}

func ScrapingStockPageStart(industryLinkList []IndustryLink) ([]StockPage, []error) {
	industryChunkList := lib.ListChunk(industryLinkList, ThreadCount)

	var allList []StockPage
	var allErrList []error

	ch1 := make(chan StockCodeLinkContainer)
	ch2 := make(chan StockPageContainer)

	go GetAllStockCodeLinkList(ch1, industryChunkList)
	go GetAllStockPage(ch1, ch2)

	for result := range ch2 {
		if result.errList != nil {
			allErrList = append(allErrList, result.errList...)
		}
		allList = append(allList, result.stockPageList...)
	}

	fmt.Println(fmt.Sprintf("len(allList) = %d", len(allList)))
	return allList, allErrList
}

func ScrapingStockCodeLinkList(indLink IndustryLink, threadCnt int) ([]StockCodeLink, error) {
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
		if s.Text() == "??????20???" {
			nextUrl, _ := s.Attr("href")
			var next *IndustryLink = &IndustryLink{}
			*next = indLink
			next.Url = fmt.Sprintf("%s%s", YahooProfileUrl, nextUrl)

			nextScList, ee := ScrapingStockCodeLinkList(*next, threadCnt)
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

func ScrapingStockPage(stockCodeLink StockCodeLink, threadCnt int) ([]StockPage, error) {
	doc, err := lib.GetPageDocument(fmt.Sprintf("%s%s", StockPageUrl, stockCodeLink.Code))
	stockPage := NewStockPage()
	stockPage.Name = stockCodeLink.Name
	stockPage.Code = stockCodeLink.Code
	fmt.Println(fmt.Sprintf("[%d] %s", threadCnt, stockCodeLink.Code))

	doc.Find("div.yjSt.info+table > tbody table tr").Each(func(_ int, tr *goquery.Selection) {
		var row = struct {
			cap string
			val string
		}{}
		tr.Find("td").Each(func(i int, td *goquery.Selection) {
			switch i {
			case 0:
				row.cap = strings.Trim(td.Text(), " ")
			case 1:
				row.val = strings.Trim(td.Text(), " ")
			}
		})

		switch row.cap {
		case "??????":
			stockPage.Topic = row.val
		case "???????????????":
			stockPage.Address = row.val
		case "????????????":
			stockPage.Tel = row.val
		case "????????????":
			stockPage.Unit = row.val
		case "?????????":
			stockPage.Market = row.val
		}
	})
	return []StockPage{*stockPage}, err
}
