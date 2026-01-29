package basic_test

import (
	"fmt"
	"log"

	tushare "github.com/fletcherlau/go-tushare"
	"github.com/fletcherlau/go-tushare/stock/basic"
)

func ExampleStockBasic() {
	// 创建客户端
	client := tushare.NewClient("your_token")

	// 获取深交所上市股票基础信息
	items, err := basic.StockBasic(client, &basic.StockBasicParams{
		Exchange:   basic.ExchangeSZSE,
		ListStatus: basic.ListStatusListed,
		Fields: []string{
			basic.StockBasicFieldTSCode,
			basic.StockBasicFieldName,
			basic.StockBasicFieldArea,
			basic.StockBasicFieldIndustry,
			basic.StockBasicFieldListDate,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("获取 %d 条记录\n", len(items))
	if len(items) > 0 {
		fmt.Printf("第一条: %s %s\n", items[0].TSCode, items[0].Name)
	}
}

func ExampleTradeCal() {
	// 创建客户端
	client := tushare.NewClient("your_token")

	// 获取上交所2024年1月的交易日历
	items, err := basic.TradeCal(client, &basic.TradeCalParams{
		Exchange:  basic.TradeCalExchangeSSE,
		StartDate: "20240101",
		EndDate:   "20240131",
		Fields: []string{
			basic.TradeCalFieldExchange,
			basic.TradeCalFieldCalDate,
			basic.TradeCalFieldIsOpen,
			basic.TradeCalFieldPretradeDate,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("获取 %d 条记录\n", len(items))
	if len(items) > 0 {
		fmt.Printf("第一条: %s 日期=%s 是否交易=%s\n",
			items[0].Exchange, items[0].CalDate, items[0].IsOpen)
	}
}
