# go-tushare

Tushare Pro HTTP API 的 Go 语言 SDK，使用 `backoff` 库实现优雅的指数退避重试机制。

## 特性

- ✅ **自动分页** - 自动处理分页，一次性获取所有数据，业务层无需关心分页逻辑
- ✅ **指数退保重试** - 基于 `github.com/cenkalti/backoff/v4` 实现智能退避策略
- ✅ **灵活重试配置** - 支持指数退避和固定间隔两种策略
- ✅ **上下文支持** - 支持 `context.Context`，可进行超时控制
- ✅ **多种配置方式** - 支持选项模式或配置结构体创建客户端
- ✅ **DataFrame 支持** - 提供类似 pandas 的 DataFrame 数据操作

## 安装

```bash
go get github.com/fletcherlau/go-tushare
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
    tushare "github.com/fletcherlau/go-tushare"
)

func main() {
    // 创建客户端（默认使用指数退避重试）
    client := tushare.NewClient("your_token_here")

    // 获取股票基础信息（自动处理分页和重试）
    resp, err := client.StockBasic(&tushare.StockBasicParams{
        Exchange:   "SZSE",
        ListStatus: "L",
        Fields:     "ts_code,name,area,industry",
    })
    if err != nil {
        log.Fatal(err)
    }

    // 打印结果
    fmt.Printf("共获取 %d 条记录\n", len(resp.Data.Items))
    for _, item := range resp.Data.Items {
        fmt.Println(item)
    }
}
```

## 客户端配置

### 方式 1：使用选项创建（推荐）

```go
// 使用指数退避（推荐）
client := tushare.NewClient("your_token",
    tushare.WithHTTPURL("https://api.tushare.pro"),
    tushare.WithTimeout(60 * time.Second),
    tushare.WithLimit(5000),                     // 每页数据条数
    tushare.WithRetries(5),                      // 最大重试次数
    tushare.WithRetryInterval(500 * time.Millisecond), // 初始重试间隔
    tushare.WithMaxInterval(30 * time.Second),   // 最大重试间隔
    tushare.WithBackoff(true),                   // 启用指数退避
)

// 使用固定间隔
client := tushare.NewClient("your_token",
    tushare.WithRetries(3),
    tushare.WithRetryInterval(2 * time.Second),  // 固定间隔 2s
    tushare.WithBackoff(false),                  // 禁用指数退避
)
```

### 方式 2：使用配置结构体

```go
conf := tushare.ClientConf{
    Token:       "your_token",
    Endpoint:    "https://api.tushare.pro",
    Limit:       5000,              // 每页数据条数
    Retries:     5,                 // 最大重试次数
    Interval:    1 * time.Second,   // 初始重试间隔
    MaxInterval: 30 * time.Second,  // 最大重试间隔
    Timeout:     30 * time.Second,  // HTTP 超时
    UseBackoff:  true,              // 使用指数退避
}

client := tushare.NewClientWithConf(conf)
```

### 方式 3：使用便捷重试配置

```go
// 默认配置（指数退避）
retryConf := tushare.DefaultRetryConfig()

// 激进配置（适合不稳定网络）
retryConf := tushare.AggressiveRetryConfig()

// 禁用重试
retryConf := tushare.NoRetryConfig()

// 使用配置创建客户端
conf := tushare.ClientConfWithRetry("your_token", retryConf)
client := tushare.NewClientWithConf(conf)
```

### 默认配置

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| Endpoint | `https://api.tushare.pro` | API 地址 |
| Limit | 5000 | 每页数据条数 |
| Retries | 3 | 最大重试次数 |
| Interval | 1s | 初始重试间隔 |
| MaxInterval | 30s | 最大重试间隔 |
| Timeout | 30s | HTTP 超时 |
| UseBackoff | true | 使用指数退避 |

## 重试机制

### 指数退避 vs 固定间隔

**指数退避（推荐）**：
- 第 1 次重试：等待 1s
- 第 2 次重试：等待 2s
- 第 3 次重试：等待 4s
- 第 4 次重试：等待 8s（但不超过 MaxInterval）

优点：避免雪崩效应，给服务器恢复时间

**固定间隔**：
- 每次重试都等待固定时间（如 2s）

### 可重试的错误

SDK 会自动对以下情况进行重试：
- 网络超时错误
- 连接错误
- 限频错误（code=40203）

不可重试的错误（立即返回）：
- API 业务错误（如权限不足、参数错误）
- HTTP 4xx 错误

### 自定义重试逻辑

```go
// 使用通用重试工具函数
err := tushare.ExecuteWithRetry(
    context.Background(),
    func() error {
        // 你的操作
        return doSomething()
    },
    5,                      // 最大重试 5 次
    true,                   // 使用指数退避
    100*time.Millisecond,   // 初始间隔
    5*time.Second,          // 最大间隔
)

// 使用带通知的重试（可记录日志或 metrics）
err := tushare.ExecuteWithRetryNotify(
    ctx,
    operation,
    5, true, 100*time.Millisecond, 5*time.Second,
    func(err error, duration time.Duration) {
        log.Printf("将在 %v 后重试，错误: %v", duration, err)
    },
)

// 标记不可重试的错误
if isPermanentError(err) {
    return tushare.PermanentError(err)
}
```

## 核心功能

### 自动分页

所有 API 方法都支持自动分页，一次性获取所有数据：

```go
// 自动获取所有分页数据（自动重试每个失败的请求）
resp, err := client.StockBasic(&tushare.StockBasicParams{
    Exchange: "SZSE",
})
// resp.Data.Items 包含所有数据，无需手动处理分页
```

如果需要单次查询（不自动分页），可以使用 `QueryOne`：

```go
// 单次查询，只返回一页数据
resp, err := client.QueryOne("stock_basic", params, fields)
```

### 超时控制

使用 context 进行超时控制：

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

resp, err := client.Query("stock_basic", params, fields, tushare.WithContext(ctx))
if err == context.DeadlineExceeded {
    fmt.Println("查询超时")
}
```

## API 说明

### 通用查询接口

```go
// 自动分页查询（带重试）
resp, err := client.Query("api_name", map[string]interface{}{
    "param1": "value1",
    "param2": "value2",
}, "field1,field2,field3")

// 单次查询（不自动分页）
resp, err := client.QueryOne("api_name", params, fields)
```

### DataFrame 操作

```go
// 获取 DataFrame（自动分页）
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
    // 检查是否为 API 错误
    if apiErr, ok := err.(*tushare.APIError); ok {
        fmt.Printf("API 错误: code=%d, msg=%s\n", apiErr.Code, apiErr.Msg)
    } else if err == context.DeadlineExceeded {
        fmt.Println("请求超时")
    } else if tushare.IsPermanentError(err) {
        fmt.Println("不可重试的错误")
    } else {
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
    Fields  []string        // 字段名列表
    Items   [][]interface{} // 数据内容
    HasMore bool            // 是否还有更多数据（内部使用）
}
```

## HTTP API 参考

本 SDK 基于 Tushare Pro HTTP API 开发，完整接口文档请参考：
https://tushare.pro/document/1?doc_id=130

## 依赖

- [github.com/cenkalti/backoff/v4](https://github.com/cenkalti/backoff) - 指数退避重试库

## License

MIT License
