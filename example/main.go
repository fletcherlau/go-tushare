package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tushare "github.com/yourusername/go-tushare"
)

func main() {
	// 从环境变量获取 token，或者直接在代码中设置
	token := os.Getenv("TUSHARE_TOKEN")
	if token == "" {
		// 请替换为您的实际 token
		token = "your_tushare_token_here"
	}

	// ========== 方式 1: 简单创建客户端 ==========
	fmt.Println("=== 方式 1: 简单创建客户端 ===")
	client := tushare.NewClient(token)
	fmt.Println("客户端创建成功（使用默认配置）")

	// ========== 方式 2: 使用自定义配置创建客户端 ==========
	fmt.Println("\n=== 方式 2: 使用自定义配置创建客户端 ===")
	conf := tushare.ClientConf{
		Token:    token,
		Endpoint: "https://api.tushare.pro",
		Limit:    3000,              // 每页 3000 条
		Retries:  5,                 // 重试 5 次
		Interval: 5 * time.Second,   // 重试间隔 5 秒
		Timeout:  30 * time.Second,  // HTTP 超时 30 秒
	}
	clientWithConf := tushare.NewClientWithConf(conf)
	fmt.Printf("配置客户端 - Limit: %d, Retries: %d, Interval: %v\n",
		conf.Limit, conf.Retries, conf.Interval)

	// ========== 方式 3: 使用选项创建客户端 ==========
	fmt.Println("\n=== 方式 3: 使用选项创建客户端 ===")
	clientWithOpts := tushare.NewClient(token,
		tushare.WithHTTPURL("https://api.tushare.pro"),
		tushare.WithLimit(5000),
		tushare.WithRetries(3),
		tushare.WithRetryInterval(10*time.Second),
		tushare.WithTimeout(60*time.Second),
	)
	fmt.Println("客户端创建成功（使用选项）")

	// 使用默认客户端进行演示
	_ = clientWithConf
	_ = clientWithOpts

	// ========== 示例 1: 获取股票基础信息（自动分页） ==========
	fmt.Println("\n=== 示例 1: 获取股票基础信息（自动分页获取所有数据） ===")
	stockBasicParams := &tushare.StockBasicParams{
		Exchange:   "SZSE",
		ListStatus: "L",
		Fields:     "ts_code,name,area,industry,list_date",
	}

	resp, err := client.StockBasic(stockBasicParams)
	if err != nil {
		log.Printf("获取股票基础信息失败: %v\n", err)
	} else {
		fmt.Printf("共获取 %d 条记录（已自动处理分页）\n", len(resp.Data.Items))
		// 打印前 5 条记录
		for i, item := range resp.Data.Items {
			if i >= 5 {
				break
			}
			fmt.Printf("记录 %d: %v\n", i+1, item)
		}
	}

	// ========== 示例 2: 获取日线行情 ==========
	fmt.Println("\n=== 示例 2: 获取日线行情 ===")
	dailyParams := &tushare.DailyParams{
		TSCode:    "000001.SZ",
		StartDate: "20240101",
		EndDate:   "20240110",
		Fields:    "ts_code,trade_date,open,high,low,close,vol",
	}

	resp, err = client.Daily(dailyParams)
	if err != nil {
		log.Printf("获取日线行情失败: %v\n", err)
	} else {
		fmt.Printf("共获取 %d 条记录\n", len(resp.Data.Items))
		df := tushare.NewDataFrame(resp)
		for i := 0; i < df.Len() && i < 5; i++ {
			fmt.Printf("日期: %s, 开盘: %.2f, 最高: %.2f, 最低: %.2f, 收盘: %.2f\n",
				df.GetString(i, "trade_date"),
				df.GetFloat64(i, "open"),
				df.GetFloat64(i, "high"),
				df.GetFloat64(i, "low"),
				df.GetFloat64(i, "close"),
			)
		}
	}

	// ========== 示例 3: 使用 DataFrame 方式查询 ==========
	fmt.Println("\n=== 示例 3: 使用 DataFrame 方式查询 ===")
	df, err := client.QueryAsDataFrame("stock_basic", map[string]interface{}{
		"exchange":    "SSE",
		"list_status": "L",
	}, "ts_code,name,area,industry")
	if err != nil {
		log.Printf("查询失败: %v\n", err)
	} else {
		fmt.Printf("共获取 %d 条记录\n", df.Len())
		for i := 0; i < df.Len() && i < 3; i++ {
			fmt.Printf("股票代码: %s, 名称: %s, 地区: %s, 行业: %s\n",
				df.GetString(i, "ts_code"),
				df.GetString(i, "name"),
				df.GetString(i, "area"),
				df.GetString(i, "industry"),
			)
		}
	}

	// ========== 示例 4: 带超时的查询 ==========
	fmt.Println("\n=== 示例 4: 带超时的查询 ===")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err = client.Query("stock_basic", map[string]interface{}{
		"exchange": "SZSE",
	}, "ts_code,name", tushare.WithContext(ctx))

	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("查询超时")
		} else {
			log.Printf("查询失败: %v\n", err)
		}
	} else {
		fmt.Printf("在超时前获取 %d 条记录\n", len(resp.Data.Items))
	}

	// ========== 示例 5: 获取指数行情 ==========
	fmt.Println("\n=== 示例 5: 获取指数行情 ===")
	indexParams := &tushare.IndexDailyParams{
		TSCode:    "000001.SH",
		StartDate: "20240101",
		EndDate:   "20240110",
	}

	resp, err = client.IndexDaily(indexParams)
	if err != nil {
		log.Printf("获取指数行情失败: %v\n", err)
	} else {
		fmt.Printf("共获取 %d 条记录\n", len(resp.Data.Items))
		df := tushare.NewDataFrame(resp)
		for i := 0; i < df.Len() && i < 5; i++ {
			fmt.Printf("日期: %s, 收盘: %.2f\n",
				df.GetString(i, "trade_date"),
				df.GetFloat64(i, "close"),
			)
		}
	}

	// ========== 示例 6: 通用查询接口 ==========
	fmt.Println("\n=== 示例 6: 通用查询接口 ===")
	resp, err = client.Query("trade_cal", map[string]interface{}{
		"exchange":   "SSE",
		"start_date": "20240101",
		"end_date":   "20240110",
		"is_open":    "1",
	}, "")
	if err != nil {
		log.Printf("获取交易日历失败: %v\n", err)
	} else {
		fmt.Printf("共获取 %d 个交易日\n", len(resp.Data.Items))
		for i, item := range resp.Data.Items {
			if i >= 5 {
				break
			}
			fmt.Printf("交易日: %v\n", item)
		}
	}

	// ========== 示例 7: 错误处理 ==========
	fmt.Println("\n=== 示例 7: 错误处理 ===")
	// 使用错误的 token 创建客户端
	badClient := tushare.NewClient("invalid_token")
	resp, err = badClient.StockBasic(&tushare.StockBasicParams{})
	if err != nil {
		if apiErr, ok := err.(*tushare.APIError); ok {
			fmt.Printf("API 错误 - 代码: %d, 消息: %s\n", apiErr.Code, apiErr.Msg)
		} else {
			fmt.Printf("请求错误: %v\n", err)
		}
	} else {
		fmt.Println("请求成功")
	}

	// ========== 示例 8: 单次查询（不处理分页） ==========
	fmt.Println("\n=== 示例 8: 单次查询（不处理分页） ===")
	// 使用 QueryOne 方法，只获取一页数据，不自动获取后续分页
	// 适用于确定数据量小的场景
	resp, err = client.QueryOne("stock_basic", map[string]interface{}{
		"ts_code": "000001.SZ",
	}, "ts_code,name")
	if err != nil {
		log.Printf("查询失败: %v\n", err)
	} else {
		fmt.Printf("单次查询获取 %d 条记录\n", len(resp.Data.Items))
	}
}
