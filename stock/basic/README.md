# stock/basic - 股票基础数据

本包提供 Tushare 股票基础数据相关接口。

## 文档参考

https://tushare.pro/document/2?doc_id=25

## 包含接口

| 接口名 | 函数 | 说明 |
|--------|------|------|
| stock_basic | `StockBasic` | 股票基础信息 |
| stock_company | `StockCompany` | 上市公司基本信息 |
| trade_cal | `TradeCal` | 交易日历 |
| namechange | `NameChange` | 股票曾用名 |
| hs_const | `HSConst` | 沪深股通成份股 |
| suspend | `StockSuspend` | 停牌信息 |
| stk_limit | `StkLimit` | 个股涨跌停 |
| stk_reward | `StkReward` | 股票质押 |

## 使用示例

```go
package main

import (
    "fmt"
    "log"
    
    tushare "github.com/yourusername/go-tushare"
    "github.com/yourusername/go-tushare/stock/basic"
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

### StockCompanyParams

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| TSCode | string | N | TS股票代码 |
| Exchange | string | N | 交易所代码 |
| Fields | string | N | 返回字段 |

### TradeCalParams

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Exchange | string | Y | 交易所代码 |
| StartDate | string | Y | 开始日期（YYYYMMDD） |
| EndDate | string | Y | 结束日期（YYYYMMDD） |
| IsOpen | string | N | 是否交易：0休市 1交易 |
| Fields | string | N | 返回字段 |

### HSConstParams

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| HsType | string | Y | 类型：SH沪股通 SZ深股通 |
| IsNew | string | N | 是否最新：1是 0否（默认1） |
| Fields | string | N | 返回字段 |
