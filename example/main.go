package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tushare "github.com/yourusername/go-tushare"
	"github.com/yourusername/go-tushare/stock/basic"
)

func main() {
	// 从环境变量获取 token，或者直接在代码中设置
	token := os.Getenv("TUSHARE_TOKEN")
	if token == "" {
		// 请替换为您的实际 token
		token = "your_tushare_token_here"
	}

	// ========== 方式 1: 简单创建客户端（使用默认指数退避）==========
	fmt.Println("=== 方式 1: 简单创建客户端（使用默认指数退避）===")
	client := tushare.NewClient(token)
	fmt.Println("客户端创建成功（使用默认配置）")

	// ========== 方式 2: 使用指数退避策略 ==========
	fmt.Println("\n=== 方式 2: 使用指数退避策略 ===")
	clientWithBackoff := tushare.NewClient(token,
		tushare.WithRetries(5),
		tushare.WithRetryInterval(500*time.Millisecond), // 初始间隔 500ms
		tushare.WithMaxInterval(30*time.Second),         // 最大间隔 30s
		tushare.WithBackoff(true),                       // 启用指数退避
	)
	fmt.Printf("指数退避配置 - 初始间隔: 500ms, 最大间隔: 30s\n")

	// ========== 方式 3: 使用固定间隔策略 ==========
	fmt.Println("\n=== 方式 3: 使用固定间隔策略 ===")
	clientWithFixed := tushare.NewClient(token,
		tushare.WithRetries(3),
		tushare.WithRetryInterval(2*time.Second), // 固定间隔 2s
		tushare.WithBackoff(false),               // 禁用指数退避
	)
	fmt.Println("固定间隔配置 - 每次重试间隔 2s")

	// ========== 方式 4: 使用配置结构体 ==========
	fmt.Println("\n=== 方式 4: 使用配置结构体 ===")
	conf := tushare.ClientConf{
		Token:       token,
		Endpoint:    "https://api.tushare.pro",
		Limit:       5000,
		Retries:     5,
		Interval:    1 * time.Second,
		MaxInterval: 60 * time.Second,
		Timeout:     30 * time.Second,
		UseBackoff:  true,
	}
	clientWithConf := tushare.NewClientWithConf(conf)
	fmt.Printf("配置客户端 - 指数退避: %v\n", conf.UseBackoff)

	// ========== 方式 5: 使用便捷重试配置 ==========
	fmt.Println("\n=== 方式 5: 使用便捷重试配置 ===")

	// 默认重试配置
	defaultRetry := tushare.DefaultRetryConfig()
	fmt.Printf("默认配置: 重试%d次, 初始间隔%v, 最大间隔%v\n",
		defaultRetry.MaxRetries, defaultRetry.InitialDelay, defaultRetry.MaxDelay)

	// 激进重试配置（适合不稳定网络）
	aggRetry := tushare.AggressiveRetryConfig()
	fmt.Printf("激进配置: 重试%d次, 初始间隔%v, 最大间隔%v\n",
		aggRetry.MaxRetries, aggRetry.InitialDelay, aggRetry.MaxDelay)

	// 禁用重试
	noRetry := tushare.NoRetryConfig()
	fmt.Printf("禁用重试: 重试%d次\n", noRetry.MaxRetries)

	// 使用配置创建客户端
	confWithRetry := tushare.ClientConfWithRetry(token, defaultRetry)
	clientWithRetryConf := tushare.NewClientWithConf(confWithRetry)
	_ = clientWithRetryConf

	// 使用前面创建的客户端进行演示
	_ = clientWithBackoff
	_ = clientWithFixed
	_ = clientWithConf

	// ========== 示例 1: 获取股票基础信息（使用 stock/basic 包）==========
	fmt.Println("\n=== 示例 1: 获取股票基础信息（使用 stock/basic 包）===")
	resp, err := basic.StockBasic(client, &basic.StockBasicParams{
		Exchange:   "SZSE",
		ListStatus: "L",
		Fields:     "ts_code,name,area,industry,list_date",
	})
	if err != nil {
		log.Printf("获取股票基础信息失败: %v\n", err)
	} else {
		fmt.Printf("共获取 %d 条记录（已自动处理分页和重试）\n", len(resp.Data.Items))
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

	// ========== 示例 6: 使用 stock/basic 包获取交易日历 ==========
	fmt.Println("\n=== 示例 6: 使用 stock/basic 包获取交易日历 ===")
	resp, err = basic.TradeCal(client, &basic.TradeCalParams{
		Exchange:  "SSE",
		StartDate: "20240101",
		EndDate:   "20240110",
		IsOpen:    "1",
		Fields:    "exchange,cal_date,is_open",
	})
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
	resp, err = basic.StockBasic(badClient, &basic.StockBasicParams{})
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

	// ========== 示例 9: 使用通用重试工具函数 ==========
	fmt.Println("\n=== 示例 9: 使用通用重试工具函数 ===")
	var attemptCount int
	err = tushare.ExecuteWithRetry(
		context.Background(),
		func() error {
			attemptCount++
			if attemptCount < 3 {
				return fmt.Errorf("模拟失败，第%d次尝试", attemptCount)
			}
			fmt.Printf("成功！共尝试 %d 次\n", attemptCount)
			return nil
		},
		5,                    // 最大重试 5 次
		true,                 // 使用指数退避
		100*time.Millisecond, // 初始间隔 100ms
		5*time.Second,        // 最大间隔 5s
	)
	if err != nil {
		log.Printf("重试最终失败: %v\n", err)
	}

	// ========== 示例 10: 使用永久错误终止重试 ==========
	fmt.Println("\n=== 示例 10: 使用永久错误终止重试 ===")
	err = tushare.ExecuteWithRetryNotify(
		context.Background(),
		func() error {
			// 某些错误不应该重试，可以直接返回永久错误
			// return tushare.PermanentError(fmt.Errorf("不应该重试的错误"))
			return nil // 这里演示成功情况
		},
		3,
		false,
		100*time.Millisecond,
		time.Second,
		func(err error, duration time.Duration) {
			fmt.Printf("即将重试，间隔: %v, 错误: %v\n", duration, err)
		},
	)
	if err != nil {
		log.Printf("执行失败: %v\n", err)
	}
}
