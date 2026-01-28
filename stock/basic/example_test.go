package basic_test

import (
	"fmt"
	"log"

	tushare "github.com/yourusername/go-tushare"
	"github.com/yourusername/go-tushare/stock/basic"
)

func ExampleStockBasic() {
	// 创建客户端
	client := tushare.NewClient("your_token")

	// 获取深交所上市股票基础信息
	resp, err := basic.StockBasic(client, &basic.StockBasicParams{
		Exchange:   "SZSE",
		ListStatus: "L",
		Fields:     "ts_code,name,area,industry,list_date",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("获取 %d 条记录\n", len(resp.Data.Items))
}

func ExampleStockCompany() {
	client := tushare.NewClient("your_token")

	// 获取上市公司基本信息
	resp, err := basic.StockCompany(client, &basic.StockCompanyParams{
		TSCode: "000001.SZ",
		Fields: "ts_code,exchange,chairman,manager,secretary",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("获取 %d 条记录\n", len(resp.Data.Items))
}

func ExampleTradeCal() {
	client := tushare.NewClient("your_token")

	// 获取2024年1月交易日历
	resp, err := basic.TradeCal(client, &basic.TradeCalParams{
		Exchange:  "SSE",
		StartDate: "20240101",
		EndDate:   "20240131",
		IsOpen:    "1", // 只获取交易日
		Fields:    "exchange,cal_date,is_open",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("2024年1月共 %d 个交易日\n", len(resp.Data.Items))
}

func ExampleHSConst() {
	client := tushare.NewClient("your_token")

	// 获取沪股通成份股
	resp, err := basic.HSConst(client, &basic.HSConstParams{
		HsType: "SH",
		IsNew:  "1",
		Fields: "ts_code,hs_type,in_date,out_date",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("沪股通成份股共 %d 条\n", len(resp.Data.Items))
}

func ExampleStockSuspend() {
	client := tushare.NewClient("your_token")

	// 获取股票停牌信息
	resp, err := basic.StockSuspend(client, &basic.StockSuspendParams{
		TSCode: "000001.SZ",
		Fields: "ts_code,suspend_date,resume_date,ann_date",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("停牌记录 %d 条\n", len(resp.Data.Items))
}
