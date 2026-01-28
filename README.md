# go-tushare

Tushare Pro HTTP API 的 Go 语言 SDK。

## 安装

```bash
go get github.com/yourusername/go-tushare
```

## 快速开始

### 1. 获取 Token

在使用前，您需要先注册 [Tushare Pro](https://tushare.pro) 账号并获取您的 API Token。

### 2. 基本使用

```go
package main

import (
    "fmt"
    "log"
    tushare "github.com/yourusername/go-tushare"
)

func main() {
    // 创建客户端
    client := tushare.NewClient("your_token_here")

    // 获取股票基础信息
    resp, err := client.StockBasic(&tushare.StockBasicParams{
        Exchange:   "SZSE",
        ListStatus: "L",
        Fields:     "ts_code,name,area,industry",
    })
    if err != nil {
        log.Fatal(err)
    }

    // 打印结果
    for _, item := range resp.Data.Items {
        fmt.Println(item)
    }
}
```

## API 说明

### 客户端创建

```go
// 基础创建
client := tushare.NewClient("your_token")

// 带选项创建
client := tushare.NewClient("your_token",
    tushare.WithHTTPURL("http://api.tushare.pro"),  // 自定义 API 地址
    tushare.WithTimeout(30),                           // 设置超时时间（秒）
)
```

### 通用查询接口

```go
resp, err := client.Query("api_name", map[string]interface{}{
    "param1": "value1",
    "param2": "value2",
}, "field1,field2,field3")
```

### DataFrame 操作

```go
// 获取 DataFrame
df, err := client.QueryAsDataFrame("stock_basic", params, fields)

// 获取行数
fmt.Println(df.Len())

// 获取字符串值
tsCode := df.GetString(0, "ts_code")

// 获取数值
closePrice := df.GetFloat64(0, "close")
```

## 支持的接口

### 股票基础数据

| 接口 | 方法 | 参数结构体 |
|------|------|-----------|
| 股票基础信息 | `StockBasic` | `StockBasicParams` |
| 上市公司信息 | `StockCompany` | `StockCompanyParams` |
| 股票曾用名 | `NameChange` | - |
| 沪深股通成份股 | `HSConst` | - |
| 停牌信息 | `StockSuspend` | - |
| 交易日历 | `TradeCal` | - |

### 行情数据

| 接口 | 方法 | 参数结构体 |
|------|------|-----------|
| 日线行情 | `Daily` | `DailyParams` |
| 周线行情 | `Weekly` | `WeeklyParams` |
| 月线行情 | `Monthly` | `MonthlyParams` |
| 每日指标 | `DailyBasic` | `DailyBasicParams` |
| 个股资金流向 | `MoneyFlow` | `MoneyFlowParams` |

### 指数数据

| 接口 | 方法 | 参数结构体 |
|------|------|-----------|
| 指数基本信息 | `IndexBasic` | `IndexBasicParams` |
| 指数日线行情 | `IndexDaily` | `IndexDailyParams` |

### 财务数据

| 接口 | 方法 | 参数结构体 |
|------|------|-----------|
| 利润表 | `Income` | `IncomeParams` |
| 资产负债表 | `BalanceSheet` | `BalanceSheetParams` |
| 现金流量表 | `CashFlow` | `CashFlowParams` |

### 期货数据

| 接口 | 方法 | 参数结构体 |
|------|------|-----------|
| 期货合约信息 | `FutBasic` | `FutBasicParams` |
| 期货日线行情 | `FutDaily` | `FutDailyParams` |

## 完整示例

查看 [example/main.go](example/main.go) 获取完整使用示例。

```bash
# 设置环境变量
export TUSHARE_TOKEN=your_token_here

# 运行示例
go run example/main.go
```

## 错误处理

```go
resp, err := client.Query("stock_basic", params, fields)
if err != nil {
    if apiErr, ok := err.(*tushare.APIError); ok {
        // API 返回的错误（如权限不足）
        fmt.Printf("API 错误: code=%d, msg=%s\n", apiErr.Code, apiErr.Msg)
    } else {
        // 网络或其他错误
        fmt.Printf("请求错误: %v\n", err)
    }
    return
}
```

## 响应数据结构

```go
type Response struct {
    Code int           // 返回码，0 表示成功
    Msg  string        // 错误信息
    Data *ResponseData // 数据
}

type ResponseData struct {
    Fields []string        // 字段名列表
    Items  [][]interface{} // 数据内容
}
```

## HTTP API 参考

本 SDK 基于 Tushare Pro HTTP API 开发，完整接口文档请参考：
https://tushare.pro/document/1?doc_id=130

## License

MIT License
