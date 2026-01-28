package tushare

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
				Fields: []string{"ts_code", "name", "area"},
				Items: [][]interface{}{
					{"000001.SZ", "平安银行", "深圳"},
					{"000002.SZ", "万科A", "深圳"},
				},
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

func TestResponse_ToRecords(t *testing.T) {
	resp := &Response{
		Code: 0,
		Data: &ResponseData{
			Fields: []string{"ts_code", "name", "price"},
			Items: [][]interface{}{
				{"000001.SZ", "平安银行", 10.5},
				{"000002.SZ", "万科A", 15.2},
			},
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

func TestClient_Options(t *testing.T) {
	// 测试自定义选项
	customURL := "http://custom.api.com"
	
	client := NewClient("test_token",
		WithHTTPURL(customURL),
		WithTimeout(60),
	)

	if client.httpURL != customURL {
		t.Errorf("期望 httpURL 为 %s，但得到 %s", customURL, client.httpURL)
	}
}
