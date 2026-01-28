// Package basic 提供 Tushare 股票基础数据接口
//
// 本包目前包含以下接口：
//   - stock_basic: 股票基础信息
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
