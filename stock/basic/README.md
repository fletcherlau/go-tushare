# stock/basic - 股票基础数据

本包提供 Tushare 股票基础数据相关接口。

## 文档参考

https://tushare.pro/document/2?doc_id=25

## 包含接口

| 接口名 | 函数 | 说明 |
|--------|------|------|
| stock_basic | `StockBasic` | 股票基础信息 |

## 使用示例

```go
package main

import (
    "fmt"
    "log"
    
    tushare "github.com/fletcherlau/go-tushare"
    "github.com/fletcherlau/go-tushare/stock/basic"
)

func main() {
    client := tushare.NewClient("your_token")
    
    // 获取股票基础信息
    resp, err := basic.StockBasic(client, &basic.StockBasicParams{
        Exchange:   "SZSE",
        ListStatus: "L",
        Fields:     "ts_code,name,area,industry",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("共获取 %d 条记录\n", len(resp.Data.Items))
}
```

## 参数说明

### StockBasicParams

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| TSCode | string | N | TS股票代码，支持多个（逗号分隔） |
| Name | string | N | 股票名称 |
| Exchange | string | N | 交易所：SSE上交所 SZSE深交所 BSE北交所 |
| Market | string | N | 市场类别：主板/创业板/科创板/CDR/北交所 |
| IsHS | string | N | 是否沪深港通：N否 H沪股通 S深股通 |
| ListStatus | string | N | 上市状态：L上市 D退市 P暂停上市 |
| Fields | string | N | 返回字段，逗号分隔 |
