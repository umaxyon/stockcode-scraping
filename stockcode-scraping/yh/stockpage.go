package yh

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"stockcode-scraping/lib"
	"strings"
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

func GetAllStockPage(c1 <-chan StockCodeLinkContainer, c2 chan<- StockPageContainer) {
	defer close(c2)
	i := 0
	for stockCodeLinkContainer := range c1 {
		chunkStockPage, errList :=
			lib.LoopScraping[StockCodeLink, StockPage](stockCodeLinkContainer.chunkStockCodeList, ScrapingStockPage, i)
		c2 <- StockPageContainer{chunkStockPage, errList}
		i++
	}
}

func ScrapingStockPageStart(industryLinkList []IndustryLink) error {
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
	return nil
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
		if s.Text() == "次の20件" {
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
		case "特色":
			stockPage.Topic = row.val
		case "本社所在地":
			stockPage.Address = row.val
		case "電話番号":
			stockPage.Tel = row.val
		case "単元株数":
			stockPage.Unit = row.val
		case "市場名":
			stockPage.Market = row.val
		}
	})
	return []StockPage{*stockPage}, err
}
