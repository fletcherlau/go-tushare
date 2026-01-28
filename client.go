package tushare

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const (
	DefaultHTTPURL       = "https://api.tushare.pro"
	DefaultTimeout       = 30 * time.Second
	DefaultLimit         = 5000
	DefaultRetries       = 3
	DefaultRetryInterval = 1 * time.Second
	DefaultMaxInterval   = 30 * time.Second

	// CodeOK 成功
	CodeOK = 0
	// CodeRateLimitExceeded 超过调用频率
	CodeRateLimitExceeded = 40203
)

// ClientConf 客户端配置
type ClientConf struct {
	Token        string        // Token
	Endpoint     string        // API 地址，默认 https://api.tushare.pro
	Limit        int           // 每页数据限制，默认 5000
	Retries      int           // 最大重试次数，默认 3
	Interval     time.Duration // 初始重试间隔，默认 1s
	MaxInterval  time.Duration // 最大重试间隔，默认 30s
	Timeout      time.Duration // HTTP 超时，默认 30s
	UseBackoff   bool          // 是否使用指数退避，默认 true
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

// WithRetries 设置最大重试次数
func WithRetries(retries int) ClientOption {
	return func(c *Client) {
		c.conf.Retries = retries
	}
}

// WithRetryInterval 设置初始重试间隔
func WithRetryInterval(interval time.Duration) ClientOption {
	return func(c *Client) {
		c.conf.Interval = interval
	}
}

// WithMaxInterval 设置最大重试间隔
func WithMaxInterval(maxInterval time.Duration) ClientOption {
	return func(c *Client) {
		c.conf.MaxInterval = maxInterval
	}
}

// WithBackoff 设置是否使用指数退避
func WithBackoff(useBackoff bool) ClientOption {
	return func(c *Client) {
		c.conf.UseBackoff = useBackoff
	}
}

// NewClient 创建新的 Tushare 客户端（兼容旧版本）
func NewClient(token string, opts ...ClientOption) *Client {
	conf := &ClientConf{
		Token:       token,
		Endpoint:    DefaultHTTPURL,
		Limit:       DefaultLimit,
		Retries:     DefaultRetries,
		Interval:    DefaultRetryInterval,
		MaxInterval: DefaultMaxInterval,
		Timeout:     DefaultTimeout,
		UseBackoff:  true,
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
	if conf.MaxInterval <= 0 {
		conf.MaxInterval = DefaultMaxInterval
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

// isRetryableError 判断错误是否可重试
func isRetryableError(resp *Response, err error) bool {
	// 网络错误可重试
	if err != nil {
		return true
	}
	// 限频错误可重试
	if resp != nil && resp.Code == CodeRateLimitExceeded {
		return true
	}
	// 服务器错误(5xx)可重试，但这里我们无法获取 HTTP 状态码
	// 其他 API 错误不重试
	return false
}

// notifyRetry 重试通知函数
func (c *Client) notifyRetry(err error, duration time.Duration) {
	// 这里可以添加日志记录或 metrics
	// 例如: log.Printf("请求失败，%v 后重试，错误: %v", duration, err)
}

// postWithRetry 发送 POST 请求（使用 backoff 实现重试）
func (c *Client) postWithRetry(apiName string, params map[string]interface{}, fields string, ctx context.Context) (*Response, error) {
	reqParams := RequestParams{
		APIName: apiName,
		Token:   c.conf.Token,
		Params:  params,
		Fields:  fields,
	}

	var resp *Response

	// 定义重试操作
	operation := func() error {
		var err error
		resp, err = c.doRequest(reqParams, ctx)

		// 成功直接返回
		if err == nil && resp.IsSuccess() {
			return nil
		}

		// 判断是否可重试
		if !isRetryableError(resp, err) {
			// 不可重试的错误，使用 backoff.Permanent 终止重试
			if err != nil {
				return backoff.Permanent(err)
			}
			// API 业务错误（非限频），不重试但返回结果
			return nil
		}

		// 可重试的错误
		if err != nil {
			return err
		}
		return fmt.Errorf("api error: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	// 配置 backoff 策略
	var b backoff.BackOff
	if c.conf.UseBackoff {
		// 指数退避
		expBackoff := backoff.NewExponentialBackOff()
		expBackoff.InitialInterval = c.conf.Interval
		expBackoff.MaxInterval = c.conf.MaxInterval
		expBackoff.Multiplier = 2
		expBackoff.RandomizationFactor = 0.1
		b = backoff.WithMaxRetries(expBackoff, uint64(c.conf.Retries))
	} else {
		// 固定间隔退避
		constBackoff := backoff.NewConstantBackOff(c.conf.Interval)
		b = backoff.WithMaxRetries(constBackoff, uint64(c.conf.Retries))
	}

	// 包装上下文支持取消
	b = backoff.WithContext(b, ctx)

	// 执行重试
	err := backoff.RetryNotify(operation, b, c.notifyRetry)
	if err != nil {
		return nil, err
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

// ==================== 便捷重试配置 ====================

// RetryConfig 是重试配置的便捷构建器
type RetryConfig struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	UseBackoff   bool
}

// DefaultRetryConfig 返回默认重试配置
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:   DefaultRetries,
		InitialDelay: DefaultRetryInterval,
		MaxDelay:     DefaultMaxInterval,
		UseBackoff:   true,
	}
}

// NoRetryConfig 返回禁用重试的配置
func NoRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries: 0,
	}
}

// AggressiveRetryConfig 返回激进重试配置（适合不稳定网络）
func AggressiveRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:   10,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     60 * time.Second,
		UseBackoff:   true,
	}
}

// ClientConfWithRetry 使用 RetryConfig 生成 ClientConf
func ClientConfWithRetry(token string, retry RetryConfig) ClientConf {
	return ClientConf{
		Token:       token,
		Endpoint:    DefaultHTTPURL,
		Limit:       DefaultLimit,
		Retries:     retry.MaxRetries,
		Interval:    retry.InitialDelay,
		MaxInterval: retry.MaxDelay,
		Timeout:     DefaultTimeout,
		UseBackoff:  retry.UseBackoff,
	}
}

// ExecuteWithRetry 使用给定的 backoff 策略执行任意函数
// 这是一个通用工具函数，可用于其他需要重试的场景
func ExecuteWithRetry(ctx context.Context, operation func() error, maxRetries int, useBackoff bool, interval, maxInterval time.Duration) error {
	var b backoff.BackOff

	if useBackoff {
		expBackoff := backoff.NewExponentialBackOff()
		expBackoff.InitialInterval = interval
		expBackoff.MaxInterval = maxInterval
		expBackoff.Multiplier = 2
		b = backoff.WithMaxRetries(expBackoff, uint64(maxRetries))
	} else {
		constBackoff := backoff.NewConstantBackOff(interval)
		b = backoff.WithMaxRetries(constBackoff, uint64(maxRetries))
	}

	b = backoff.WithContext(b, ctx)

	return backoff.Retry(operation, b)
}

// ExecuteWithRetryNotify 带通知的重试执行
func ExecuteWithRetryNotify(ctx context.Context, operation func() error, maxRetries int, useBackoff bool, interval, maxInterval time.Duration, notify func(err error, duration time.Duration)) error {
	var b backoff.BackOff

	if useBackoff {
		expBackoff := backoff.NewExponentialBackOff()
		expBackoff.InitialInterval = interval
		expBackoff.MaxInterval = maxInterval
		expBackoff.Multiplier = 2
		b = backoff.WithMaxRetries(expBackoff, uint64(maxRetries))
	} else {
		constBackoff := backoff.NewConstantBackOff(interval)
		b = backoff.WithMaxRetries(constBackoff, uint64(maxRetries))
	}

	b = backoff.WithContext(b, ctx)

	return backoff.RetryNotify(operation, b, notify)
}

// PermanentError 包装一个错误为不可重试错误
// 使用示例: if shouldNotRetry(err) { return tushare.PermanentError(err) }
func PermanentError(err error) error {
	return backoff.Permanent(err)
}

// IsPermanentError 检查错误是否为不可重试错误
func IsPermanentError(err error) bool {
	var permanent *backoff.PermanentError
	return errors.As(err, &permanent)
}
