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
