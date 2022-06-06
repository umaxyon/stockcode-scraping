package db

import (
	"stockcode-scraping/db/test"
	"stockcode-scraping/yh"
	"testing"
)

func TestMain(m *testing.M) {
	test.PrepareTestAspect(func() int {
		return m.Run()
	}, "../../dynamo_ddl.yaml")
}

func TestDynamo(t *testing.T) {
	t.Run("dynamo test", func(t *testing.T) {

		stockPageList := []yh.StockPage{
			{Name: "aaa", Code: "1111", Tel: "000-0000-0001", Address: "大阪府大阪市梅田1-1-1", Topic: "aaatopic", Market: "東証グレイス1", Unit: "101"},
			{Name: "aab", Code: "1112", Tel: "000-0000-0002", Address: "大阪府大阪市梅田1-1-2", Topic: "aabtopic", Market: "東証グレイス2", Unit: "102"},
			{Name: "aac", Code: "1113", Tel: "000-0000-0003", Address: "大阪府大阪市梅田1-1-3", Topic: "aactopic", Market: "東証グレイス3", Unit: "103"},
		}

		accessor := NewAccessor()
		accessor.SaveStCode(stockPageList)
		accessor.Query()
	})
}
