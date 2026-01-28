// Package basic 提供 Tushare 股票基础数据接口
//
// 本包包含以下接口：
//   - stock_basic: 股票基础信息
//   - stock_company: 上市公司基本信息
//   - trade_cal: 交易日历
//   - namechange: 股票曾用名
//   - hs_const: 沪深股通成份股
//   - suspend: 停牌信息
//   - stk_limit: 个股涨跌停
//   - stk_reward: 股票质押
//
// 文档参考: https://tushare.pro/document/2?doc_id=25
//
// 使用示例：
//
//	import (
//	    tushare "github.com/fletcherlau/go-tushare"
//	    "github.com/fletcherlau/go-tushare/stock/basic"
//	)
//
//	func main() {
//	    client := tushare.NewClient("your_token")
//	    
//	    // 获取股票基础信息
//	    resp, err := basic.StockBasic(client, &basic.StockBasicParams{
//	        Exchange: "SZSE",
//	        Fields:   "ts_code,name,area",
//	    })
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    
//	    fmt.Printf("获取 %d 条记录\n", len(resp.Data.Items))
//	}
//
package basic
