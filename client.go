package tushare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultHTTPURL = "http://api.tushare.pro"
	DefaultTimeout = 30 * time.Second
)

// Client Tushare API 客户端
type Client struct {
	token      string
	httpURL    string
	httpClient *http.Client
}

// ClientOption 客户端配置选项
type ClientOption func(*Client)

// WithHTTPURL 设置自定义 API 地址
func WithHTTPURL(url string) ClientOption {
	return func(c *Client) {
		c.httpURL = url
	}
}

// WithTimeout 设置 HTTP 超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithHTTPClient 设置自定义 HTTP 客户端
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// NewClient 创建新的 Tushare 客户端
func NewClient(token string, opts ...ClientOption) *Client {
	client := &Client{
		token:   token,
		httpURL: DefaultHTTPURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
	
	for _, opt := range opts {
		opt(client)
	}
	
	return client
}

// RequestParams 请求参数
type RequestParams struct {
	APIName string                 `json:"api_name"`
	Token   string                 `json:"token"`
	Params  map[string]interface{} `json:"params,omitempty"`
	Fields  string                 `json:"fields,omitempty"`
}

// Query 执行通用查询
// apiName: 接口名称，如 "stock_basic"
// params: 接口参数，如 {"list_status": "L"}
// fields: 需要返回的字段，逗号分隔，如 "ts_code,name,area"
func (c *Client) Query(apiName string, params map[string]interface{}, fields string) (*Response, error) {
	reqParams := RequestParams{
		APIName: apiName,
		Token:   c.token,
		Params:  params,
		Fields:  fields,
	}
	
	return c.doRequest(reqParams)
}

// QueryWithContext 执行带上下文的通用查询（可用于超时控制）
func (c *Client) QueryWithContext(apiName string, params map[string]interface{}, fields string) (*Response, error) {
	// 目前先简单实现，后续可以添加 context.Context 支持
	return c.Query(apiName, params, fields)
}

// doRequest 执行 HTTP 请求
func (c *Client) doRequest(reqParams RequestParams) (*Response, error) {
	jsonBody, err := json.Marshal(reqParams)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}
	
	req, err := http.NewRequest(http.MethodPost, c.httpURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: status=%d, body=%s", resp.StatusCode, string(body))
	}
	
	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w, body=%s", err, string(body))
	}
	
	if !result.IsSuccess() {
		return &result, &APIError{
			Code: result.Code,
			Msg:  result.Msg,
		}
	}
	
	return &result, nil
}

// QueryAsDataFrame 执行查询并返回 DataFrame
func (c *Client) QueryAsDataFrame(apiName string, params map[string]interface{}, fields string) (*DataFrame, error) {
	resp, err := c.Query(apiName, params, fields)
	if err != nil {
		return nil, err
	}
	return NewDataFrame(resp), nil
}
