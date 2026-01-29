package market_test

import (
	"fmt"
	"log"

	tushare "github.com/fletcherlau/go-tushare"
	"github.com/fletcherlau/go-tushare/stock/market"
)

func ExampleDaily() {
	// 创建客户端
	client := tushare.NewClient("your_token")

	// 获取平安银行2024年1月的日线行情
	items, err := market.Daily(client, &market.DailyParams{
		TSCode:    "000001.SZ",
		StartDate: "20240101",
		EndDate:   "20240131",
		Fields: []string{
			market.DailyFieldTradeDate,
			market.DailyFieldOpen,
			market.DailyFieldHigh,
			market.DailyFieldLow,
			market.DailyFieldClose,
			market.DailyFieldVol,
			market.DailyFieldAmount,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("获取 %d 条记录\n", len(items))
	if len(items) > 0 {
		fmt.Printf("第一条: 日期=%s 开盘=%.2f 收盘=%.2f\n",
			items[0].TradeDate, items[0].Open, items[0].Close)
	}
}

func ExampleAdjFactor() {
	// 创建客户端
	client := tushare.NewClient("your_token")

	// 获取平安银行的复权因子
	items, err := market.AdjFactor(client, &market.AdjFactorParams{
		TSCode:    "000001.SZ",
		StartDate: "20240101",
		EndDate:   "20240131",
		Fields: []string{
			market.AdjFactorFieldTradeDate,
			market.AdjFactorFieldAdjFactor,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("获取 %d 条记录\n", len(items))
	if len(items) > 0 {
		fmt.Printf("第一条: 日期=%s 复权因子=%.6f\n",
			items[0].TradeDate, items[0].AdjFactor)
	}
}

func ExampleDailyBasic() {
	// 创建客户端
	client := tushare.NewClient("your_token")

	// 获取平安银行的每日指标
	items, err := market.DailyBasic(client, &market.DailyBasicParams{
		TSCode:    "000001.SZ",
		StartDate: "20240101",
		EndDate:   "20240131",
		Fields: []string{
			market.DailyBasicFieldTradeDate,
			market.DailyBasicFieldClose,
			market.DailyBasicFieldPE,
			market.DailyBasicFieldPB,
			market.DailyBasicFieldTotalMV,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("获取 %d 条记录\n", len(items))
	if len(items) > 0 {
		fmt.Printf("第一条: 日期=%s 收盘价=%.2f PE=%.2f PB=%.2f 总市值=%.0f万\n",
			items[0].TradeDate, items[0].Close, items[0].PE, items[0].PB, items[0].TotalMV)
	}
}
