package tushare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultHTTPURL     = "https://api.tushare.pro"
	DefaultTimeout     = 30 * time.Second
	DefaultLimit       = 5000
	DefaultRetries     = 3
	DefaultRetryInterval = 10 * time.Second
	
	// CodeOK 成功
	CodeOK = 0
	// CodeRateLimitExceeded 超过调用频率
	CodeRateLimitExceeded = 40203
)

// ClientConf 客户端配置
type ClientConf struct {
	Token    string        // Token
	Endpoint string        // API 地址，默认 https://api.tushare.pro
	Limit    int           // 每页数据限制，默认 5000
	Retries  int           // 重试次数，默认 3
	Interval time.Duration // 重试间隔，默认 10s
	Timeout  time.Duration // HTTP 超时，默认 30s
}

// Client Tushare API 客户端
type Client struct {
	conf   *ClientConf
	client *http.Client
}

// ClientOption 客户端配置选项
type ClientOption func(*Client)

// WithHTTPURL 设置自定义 API 地址
func WithHTTPURL(url string) ClientOption {
	return func(c *Client) {
		c.conf.Endpoint = url
	}
}

// WithTimeout 设置 HTTP 超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.conf.Timeout = timeout
		c.client.Timeout = timeout
	}
}

// WithHTTPClient 设置自定义 HTTP 客户端
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.client = client
	}
}

// WithLimit 设置分页大小
func WithLimit(limit int) ClientOption {
	return func(c *Client) {
		c.conf.Limit = limit
	}
}

// WithRetries 设置重试次数
func WithRetries(retries int) ClientOption {
	return func(c *Client) {
		c.conf.Retries = retries
	}
}

// WithRetryInterval 设置重试间隔
func WithRetryInterval(interval time.Duration) ClientOption {
	return func(c *Client) {
		c.conf.Interval = interval
	}
}

// NewClient 创建新的 Tushare 客户端（兼容旧版本）
func NewClient(token string, opts ...ClientOption) *Client {
	conf := &ClientConf{
		Token:    token,
		Endpoint: DefaultHTTPURL,
		Limit:    DefaultLimit,
		Retries:  DefaultRetries,
		Interval: DefaultRetryInterval,
		Timeout:  DefaultTimeout,
	}

	client := &Client{
		conf: conf,
		client: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// NewClientWithConf 使用配置创建客户端
func NewClientWithConf(conf ClientConf, opts ...ClientOption) *Client {
	// 设置默认值
	if conf.Endpoint == "" {
		conf.Endpoint = DefaultHTTPURL
	}
	if conf.Limit <= 0 {
		conf.Limit = DefaultLimit
	}
	if conf.Retries <= 0 {
		conf.Retries = DefaultRetries
	}
	if conf.Interval <= 0 {
		conf.Interval = DefaultRetryInterval
	}
	if conf.Timeout <= 0 {
		conf.Timeout = DefaultTimeout
	}

	client := &Client{
		conf: &conf,
		client: &http.Client{
			Timeout: conf.Timeout,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// QueryOption 查询选项
type QueryOption func(*queryOptions)

type queryOptions struct {
	ctx context.Context
}

// WithContext 添加上下文选项（用于超时控制）
func WithContext(ctx context.Context) QueryOption {
	return func(o *queryOptions) {
		o.ctx = ctx
	}
}

// defaultQueryOptions 默认查询选项
func defaultQueryOptions() *queryOptions {
	return &queryOptions{
		ctx: context.Background(),
	}
}

// RequestParams 请求参数
type RequestParams struct {
	APIName string                 `json:"api_name"`
	Token   string                 `json:"token"`
	Params  map[string]interface{} `json:"params,omitempty"`
	Fields  string                 `json:"fields,omitempty"`
}

// Query 执行通用查询（自动处理分页，一次性获取所有数据）
func (c *Client) Query(apiName string, params map[string]interface{}, fields string, opts ...QueryOption) (*Response, error) {
	// 合并选项
	options := defaultQueryOptions()
	for _, opt := range opts {
		opt(options)
	}

	// 复制参数，避免修改原始参数
	newParams := make(map[string]interface{})
	for k, v := range params {
		newParams[k] = v
	}

	// 设置分页参数
	newParams["limit"] = c.conf.Limit
	newParams["offset"] = 0

	// 合并所有数据
	allItems := make([][]interface{}, 0)
	var respFields []string

	for {
		// 检查上下文是否已取消
		select {
		case <-options.ctx.Done():
			return nil, options.ctx.Err()
		default:
		}

		resp, err := c.postWithRetry(apiName, newParams, fields, options.ctx)
		if err != nil {
			return nil, err
		}

		if !resp.IsSuccess() {
			return resp, &APIError{
				Code: resp.Code,
				Msg:  resp.Msg,
			}
		}

		if resp.Data != nil {
			respFields = resp.Data.Fields
			allItems = append(allItems, resp.Data.Items...)

			// 如果没有更多数据，退出循环
			if !resp.Data.HasMore {
				break
			}
		} else {
			break
		}

		// 更新 offset 继续获取下一页
		offset, _ := newParams["offset"].(int)
		newParams["offset"] = offset + c.conf.Limit
	}

	// 构造合并后的响应
	return &Response{
		Code: CodeOK,
		Msg:  "",
		Data: &ResponseData{
			Fields:  respFields,
			Items:   allItems,
			HasMore: false,
		},
	}, nil
}

// QueryOne 执行单次查询（不处理分页，用于确定数据量小的场景）
func (c *Client) QueryOne(apiName string, params map[string]interface{}, fields string, opts ...QueryOption) (*Response, error) {
	options := defaultQueryOptions()
	for _, opt := range opts {
		opt(options)
	}

	return c.postWithRetry(apiName, params, fields, options.ctx)
}

// postWithRetry 发送 POST 请求（带重试机制）
func (c *Client) postWithRetry(apiName string, params map[string]interface{}, fields string, ctx context.Context) (*Response, error) {
	reqParams := RequestParams{
		APIName: apiName,
		Token:   c.conf.Token,
		Params:  params,
		Fields:  fields,
	}

	var lastErr error
	var resp *Response

	for i := 0; i < c.conf.Retries; i++ {
		// 检查上下文
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		resp, lastErr = c.doRequest(reqParams, ctx)
		if lastErr == nil && resp.IsSuccess() {
			return resp, nil
		}

		// 如果是限频错误，等待后重试
		if resp != nil && resp.Code == CodeRateLimitExceeded {
			time.Sleep(c.conf.Interval)
			continue
		}

		// 如果是网络错误，等待后重试
		if lastErr != nil && i < c.conf.Retries-1 {
			time.Sleep(c.conf.Interval)
			continue
		}

		// 其他错误直接返回
		if lastErr == nil && !resp.IsSuccess() {
			return resp, nil
		}
		break
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return resp, nil
}

// doRequest 执行 HTTP 请求
func (c *Client) doRequest(reqParams RequestParams, ctx context.Context) (*Response, error) {
	jsonBody, err := json.Marshal(reqParams)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.conf.Endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
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

	return &result, nil
}

// QueryAsDataFrame 执行查询并返回 DataFrame（自动分页）
func (c *Client) QueryAsDataFrame(apiName string, params map[string]interface{}, fields string, opts ...QueryOption) (*DataFrame, error) {
	resp, err := c.Query(apiName, params, fields, opts...)
	if err != nil {
		return nil, err
	}
	return NewDataFrame(resp), nil
}

// ==================== 向后兼容的方法 ====================

// QueryWithContext 执行带上下文的通用查询（兼容旧版本，实际等价于 Query）
func (c *Client) QueryWithContext(apiName string, params map[string]interface{}, fields string) (*Response, error) {
	return c.Query(apiName, params, fields, WithContext(context.Background()))
}
