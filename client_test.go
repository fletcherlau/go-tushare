package tushare

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestClient_Query(t *testing.T) {
	// 创建模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		if r.Method != http.MethodPost {
			t.Errorf("期望 POST 方法，但得到 %s", r.Method)
		}

		// 验证 Content-Type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("期望 Content-Type 为 application/json，但得到 %s", contentType)
		}

		// 解析请求体
		var reqParams RequestParams
		if err := json.NewDecoder(r.Body).Decode(&reqParams); err != nil {
			t.Errorf("解析请求体失败: %v", err)
			return
		}

		// 验证请求参数
		if reqParams.APIName != "stock_basic" {
			t.Errorf("期望 api_name 为 stock_basic，但得到 %s", reqParams.APIName)
		}
		if reqParams.Token != "test_token" {
			t.Errorf("期望 token 为 test_token，但得到 %s", reqParams.Token)
		}

		// 返回模拟响应
		response := Response{
			Code: 0,
			Msg:  "",
			Data: &ResponseData{
				Fields:  []string{"ts_code", "name", "area"},
				Items: [][]interface{}{
					{"000001.SZ", "平安银行", "深圳"},
					{"000002.SZ", "万科A", "深圳"},
				},
				HasMore: false,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 创建客户端
	client := NewClient("test_token", WithHTTPURL(server.URL))

	// 测试查询
	resp, err := client.Query("stock_basic", map[string]interface{}{
		"list_status": "L",
	}, "ts_code,name,area")

	if err != nil {
		t.Errorf("查询失败: %v", err)
		return
	}

	if !resp.IsSuccess() {
		t.Errorf("期望返回成功，但 code=%d", resp.Code)
	}

	if resp.Data == nil {
		t.Fatal("响应数据为空")
	}

	if len(resp.Data.Fields) != 3 {
		t.Errorf("期望 3 个字段，但得到 %d", len(resp.Data.Fields))
	}

	if len(resp.Data.Items) != 2 {
		t.Errorf("期望 2 条记录，但得到 %d", len(resp.Data.Items))
	}
}

func TestClient_QueryPagination(t *testing.T) {
	// 请求计数器
	var requestCount int32

	// 创建模拟服务器，模拟分页数据
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)

		var reqParams RequestParams
		if err := json.NewDecoder(r.Body).Decode(&reqParams); err != nil {
			t.Errorf("解析请求体失败: %v", err)
			return
		}

		offset, _ := reqParams.Params["offset"].(float64)

		var response Response
		if offset == 0 {
			// 第一页数据
			response = Response{
				Code: 0,
				Data: &ResponseData{
					Fields: []string{"ts_code", "name"},
					Items: [][]interface{}{
						{"000001.SZ", "平安银行"},
						{"000002.SZ", "万科A"},
					},
					HasMore: true,
				},
			}
		} else {
			// 第二页数据
			response = Response{
				Code: 0,
				Data: &ResponseData{
					Fields: []string{"ts_code", "name"},
					Items: [][]interface{}{
						{"000003.SZ", "PT金田A"},
					},
					HasMore: false,
				},
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 创建客户端，设置每页 2 条
	client := NewClient("test_token",
		WithHTTPURL(server.URL),
		WithLimit(2),
	)

	// 测试查询 - 应该自动获取两页数据
	resp, err := client.Query("stock_basic", map[string]interface{}{}, "ts_code,name")

	if err != nil {
		t.Errorf("查询失败: %v", err)
		return
	}

	// 验证请求次数
	if requestCount != 2 {
		t.Errorf("期望请求 2 次（分页），但实际请求 %d 次", requestCount)
	}

	// 验证合并后的数据
	if len(resp.Data.Items) != 3 {
		t.Errorf("期望合并后 3 条记录，但得到 %d", len(resp.Data.Items))
	}
}

func TestClient_QueryWithContext(t *testing.T) {
	// 创建延迟响应的模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		response := Response{
			Code: 0,
			Data: &ResponseData{
				Fields:  []string{"ts_code"},
				Items:   [][]interface{}{{"000001.SZ"}},
				HasMore: false,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test_token", WithHTTPURL(server.URL))

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// 测试超时
	_, err := client.Query("stock_basic", map[string]interface{}{}, "", WithContext(ctx))

	if err == nil {
		t.Error("期望返回超时错误，但没有")
	}
}

func TestClient_QueryError(t *testing.T) {
	// 创建返回错误的模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			Code: 2002,
			Msg:  "没有权限",
			Data: nil,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("invalid_token", WithHTTPURL(server.URL))

	resp, err := client.Query("stock_basic", nil, "")

	// 期望返回 APIError
	if err == nil {
		t.Error("期望返回错误，但没有")
		return
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Errorf("期望错误类型为 *APIError，但得到 %T", err)
		return
	}

	if apiErr.Code != 2002 {
		t.Errorf("期望错误码为 2002，但得到 %d", apiErr.Code)
	}

	if apiErr.Msg != "没有权限" {
		t.Errorf("期望错误消息为'没有权限'，但得到 %s", apiErr.Msg)
	}

	// 验证响应对象也被返回
	if resp == nil {
		t.Error("期望返回响应对象，但没有")
	}
}

func TestClient_BackoffRetry(t *testing.T) {
	// 请求计数器
	var requestCount int32

	// 创建模拟服务器，前两次返回限频错误
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)

		if count <= 2 {
			// 返回限频错误
			response := Response{
				Code: CodeRateLimitExceeded,
				Msg:  "超过调用频率",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// 第三次返回成功
		response := Response{
			Code: 0,
			Data: &ResponseData{
				Fields:  []string{"ts_code"},
				Items:   [][]interface{}{{"000001.SZ"}},
				HasMore: false,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 创建客户端，使用指数退避
	client := NewClient("test_token",
		WithHTTPURL(server.URL),
		WithRetries(5),
		WithRetryInterval(10*time.Millisecond),
		WithBackoff(true), // 使用指数退避
	)

	// 测试查询 - 应该自动重试
	resp, err := client.Query("stock_basic", map[string]interface{}{}, "")

	if err != nil {
		t.Errorf("查询失败: %v", err)
		return
	}

	// 验证请求次数（3次：2次失败+1次成功）
	if requestCount != 3 {
		t.Errorf("期望请求 3 次（含重试），但实际请求 %d 次", requestCount)
	}

	if !resp.IsSuccess() {
		t.Error("期望最终成功")
	}
}

func TestClient_NoBackoffRetry(t *testing.T) {
	// 请求计数器
	var requestCount int32

	// 创建模拟服务器，第一次返回限频错误
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)

		if count == 1 {
			// 返回限频错误
			response := Response{
				Code: CodeRateLimitExceeded,
				Msg:  "超过调用频率",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// 第二次返回成功
		response := Response{
			Code: 0,
			Data: &ResponseData{
				Fields:  []string{"ts_code"},
				Items:   [][]interface{}{{"000001.SZ"}},
				HasMore: false,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 创建客户端，使用固定间隔
	client := NewClient("test_token",
		WithHTTPURL(server.URL),
		WithRetries(3),
		WithRetryInterval(50*time.Millisecond),
		WithBackoff(false), // 使用固定间隔
	)

	// 测试查询
	resp, err := client.Query("stock_basic", map[string]interface{}{}, "")

	if err != nil {
		t.Errorf("查询失败: %v", err)
		return
	}

	// 验证请求次数
	if requestCount != 2 {
		t.Errorf("期望请求 2 次，但实际请求 %d 次", requestCount)
	}

	if !resp.IsSuccess() {
		t.Error("期望最终成功")
	}
}

func TestClient_RetryExhausted(t *testing.T) {
	// 请求计数器
	var requestCount int32

	// 创建总是返回限频错误的模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		response := Response{
			Code: CodeRateLimitExceeded,
			Msg:  "超过调用频率",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 创建客户端，只重试 2 次
	client := NewClient("test_token",
		WithHTTPURL(server.URL),
		WithRetries(2),
		WithRetryInterval(10*time.Millisecond),
	)

	// 测试查询 - 应该最终失败
	_, err := client.Query("stock_basic", map[string]interface{}{}, "")

	if err == nil {
		t.Error("期望返回错误（重试耗尽），但没有")
		return
	}

	// 验证请求次数：原始请求 + 2次重试 = 3次
	if requestCount != 3 {
		t.Errorf("期望请求 3 次（原始+2次重试），但实际请求 %d 次", requestCount)
	}
}

func TestClient_NewClientWithConf(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			Code: 0,
			Data: &ResponseData{
				Fields:  []string{"ts_code"},
				Items:   [][]interface{}{{"000001.SZ"}},
				HasMore: false,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 使用配置创建客户端
	conf := ClientConf{
		Token:       "test_token",
		Endpoint:    server.URL,
		Limit:       100,
		Retries:     5,
		Interval:    100 * time.Millisecond,
		MaxInterval: 5 * time.Second,
		Timeout:     30 * time.Second,
		UseBackoff:  true,
	}

	client := NewClientWithConf(conf)

	if client.conf.Token != "test_token" {
		t.Errorf("期望 Token 为 test_token，但得到 %s", client.conf.Token)
	}

	if client.conf.Limit != 100 {
		t.Errorf("期望 Limit 为 100，但得到 %d", client.conf.Limit)
	}

	// 测试调用
	resp, err := client.Query("test", nil, "")
	if err != nil {
		t.Errorf("查询失败: %v", err)
	}
	if !resp.IsSuccess() {
		t.Error("期望成功")
	}
}

func TestClient_RetryConfig(t *testing.T) {
	// 测试默认重试配置
	defaultConfig := DefaultRetryConfig()
	if defaultConfig.MaxRetries != DefaultRetries {
		t.Errorf("期望默认重试次数 %d，但得到 %d", DefaultRetries, defaultConfig.MaxRetries)
	}

	// 测试禁用重试配置
	noRetryConfig := NoRetryConfig()
	if noRetryConfig.MaxRetries != 0 {
		t.Errorf("期望禁用重试，但得到重试次数 %d", noRetryConfig.MaxRetries)
	}

	// 测试激进重试配置
	aggConfig := AggressiveRetryConfig()
	if aggConfig.MaxRetries != 10 {
		t.Errorf("期望激进重试次数 10，但得到 %d", aggConfig.MaxRetries)
	}

	// 测试使用 RetryConfig 创建 ClientConf
	conf := ClientConfWithRetry("test_token", defaultConfig)
	if conf.Token != "test_token" {
		t.Errorf("期望 Token 为 test_token，但得到 %s", conf.Token)
	}
	if conf.Retries != defaultConfig.MaxRetries {
		t.Errorf("期望 Retries 为 %d，但得到 %d", defaultConfig.MaxRetries, conf.Retries)
	}
}

func TestResponse_ToRecords(t *testing.T) {
	resp := &Response{
		Code: 0,
		Data: &ResponseData{
			Fields: []string{"ts_code", "name", "price"},
			Items: [][]interface{}{
				{"000001.SZ", "平安银行", 10.5},
				{"000002.SZ", "万科A", 15.2},
			},
			HasMore: false,
		},
	}

	records := resp.ToRecords()

	if len(records) != 2 {
		t.Errorf("期望 2 条记录，但得到 %d", len(records))
	}

	// 验证第一条记录
	if records[0]["ts_code"] != "000001.SZ" {
		t.Errorf("期望 ts_code 为 000001.SZ，但得到 %v", records[0]["ts_code"])
	}
	if records[0]["name"] != "平安银行" {
		t.Errorf("期望 name 为 平安银行，但得到 %v", records[0]["name"])
	}
	if records[0]["price"] != 10.5 {
		t.Errorf("期望 price 为 10.5，但得到 %v", records[0]["price"])
	}
}

func TestDataFrame(t *testing.T) {
	resp := &Response{
		Code: 0,
		Data: &ResponseData{
			Fields: []string{"ts_code", "name", "open", "close"},
			Items: [][]interface{}{
				{"000001.SZ", "平安银行", 10.0, 10.5},
				{"000002.SZ", "万科A", 15.0, 15.2},
			},
			HasMore: false,
		},
	}

	df := NewDataFrame(resp)

	if df.Len() != 2 {
		t.Errorf("期望 DataFrame 长度为 2，但得到 %d", df.Len())
	}

	// 测试 GetString
	if df.GetString(0, "ts_code") != "000001.SZ" {
		t.Errorf("期望 ts_code 为 000001.SZ，但得到 %s", df.GetString(0, "ts_code"))
	}

	// 测试 GetFloat64
	if df.GetFloat64(0, "open") != 10.0 {
		t.Errorf("期望 open 为 10.0，但得到 %f", df.GetFloat64(0, "open"))
	}

	if df.GetFloat64(0, "close") != 10.5 {
		t.Errorf("期望 close 为 10.5，但得到 %f", df.GetFloat64(0, "close"))
	}

	// 测试不存在的列
	if df.GetString(0, "nonexistent") != "" {
		t.Error("不存在的列应该返回空字符串")
	}

	// 测试越界
	if df.GetString(10, "ts_code") != "" {
		t.Error("越界访问应该返回空字符串")
	}
}

func TestClient_StockBasic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqParams RequestParams
		if err := json.NewDecoder(r.Body).Decode(&reqParams); err != nil {
			t.Errorf("解析请求体失败: %v", err)
			return
		}

		// 验证参数
		if reqParams.APIName != "stock_basic" {
			t.Errorf("期望 api_name 为 stock_basic，但得到 %s", reqParams.APIName)
		}

		// 验证 list_status 默认值
		if reqParams.Params["list_status"] != "L" {
			t.Errorf("期望 list_status 默认为 L，但得到 %v", reqParams.Params["list_status"])
		}

		response := Response{
			Code: 0,
			Data: &ResponseData{
				Fields:  []string{"ts_code", "name"},
				Items:   [][]interface{}{{"000001.SZ", "平安银行"}},
				HasMore: false,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test_token", WithHTTPURL(server.URL))

	params := &StockBasicParams{
		Exchange: "SZSE",
		Fields:   "ts_code,name",
	}

	resp, err := client.StockBasic(params)
	if err != nil {
		t.Errorf("查询失败: %v", err)
		return
	}

	if !resp.IsSuccess() {
		t.Errorf("期望成功，但 code=%d", resp.Code)
	}
}

func TestExecuteWithRetry(t *testing.T) {
	// 测试通用重试函数
	var count int
	operation := func() error {
		count++
		if count < 3 {
			return fmt.Errorf("temporary error")
		}
		return nil
	}

	ctx := context.Background()
	err := ExecuteWithRetry(ctx, operation, 5, false, 10*time.Millisecond, time.Second)

	if err != nil {
		t.Errorf("期望成功，但得到错误: %v", err)
	}

	if count != 3 {
		t.Errorf("期望执行 3 次，但实际执行 %d 次", count)
	}
}

func TestPermanentError(t *testing.T) {
	// 测试永久错误
	err := PermanentError(fmt.Errorf("permanent error"))

	if !IsPermanentError(err) {
		t.Error("期望是永久错误")
	}

	// 测试普通错误
	normalErr := fmt.Errorf("normal error")
	if IsPermanentError(normalErr) {
		t.Error("普通错误不应是永久错误")
	}
}
