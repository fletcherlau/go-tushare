// Package financial 提供 Tushare 股票财务数据接口
//
// 本包目前包含以下接口：
//   - income: 利润表
//   - balancesheet: 资产负债表
//   - cashflow: 现金流量表
//   - fina_indicator: 财务指标数据
//
// 文档参考:
//   - income: https://tushare.pro/document/2?doc_id=33
//   - balancesheet: https://tushare.pro/document/2?doc_id=36
//   - cashflow: https://tushare.pro/document/2?doc_id=44
//   - fina_indicator: https://tushare.pro/document/2?doc_id=79
//
// 使用示例：
//
//	import (
//	    tushare "github.com/fletcherlau/go-tushare"
//	    "github.com/fletcherlau/go-tushare/stock/financial"
//	)
//
//	func main() {
//	    client := tushare.NewClient("your_token")
//
//	    // 获取利润表
//	    items, err := financial.Income(client, &financial.IncomeParams{
//	        TSCode: "000001.SZ",
//	    })
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    fmt.Printf("获取 %d 条记录\n", len(items))
//	}
//
package financial
