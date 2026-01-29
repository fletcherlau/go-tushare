// Package market 提供 Tushare 股票行情数据接口
//
// 本包目前包含以下接口：
//   - daily: A股日线行情
//   - adj_factor: 复权因子
//   - daily_basic: 每日指标
//
// 文档参考:
//   - daily: https://tushare.pro/document/2?doc_id=27
//   - adj_factor: https://tushare.pro/document/2?doc_id=28
//   - daily_basic: https://tushare.pro/document/2?doc_id=32
//
// 使用示例：
//
//	import (
//	    tushare "github.com/fletcherlau/go-tushare"
//	    "github.com/fletcherlau/go-tushare/stock/market"
//	)
//
//	func main() {
//	    client := tushare.NewClient("your_token")
//
//	    // 获取日线行情
//	    items, err := market.Daily(client, &market.DailyParams{
//	        TSCode: "000001.SZ",
//	        Fields: []string{market.DailyFieldTradeDate, market.DailyFieldClose},
//	    })
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    fmt.Printf("获取 %d 条记录\n", len(items))
//	}
//
package market
