package tushare

import (
	"fmt"
)

// ==================== 日线行情 ====================

// DailyParams 日线行情参数
type DailyParams struct {
	TSCode    string // 股票代码（支持多个股票同时提取，逗号分隔）
	TradeDate string // 交易日期（YYYYMMDD）
	StartDate string // 开始日期(YYYYMMDD)
	EndDate   string // 结束日期(YYYYMMDD)
	Fields    string // 返回字段
}

// Daily 获取日线行情（自动处理分页）
func (c *Client) Daily(params *DailyParams, opts ...QueryOption) (*Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.TradeDate != "" {
		reqParams["trade_date"] = params.TradeDate
	}
	if params.StartDate != "" {
		reqParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		reqParams["end_date"] = params.EndDate
	}

	return c.Query("daily", reqParams, params.Fields, opts...)
}

// ==================== 周线行情 ====================

// WeeklyParams 周线行情参数
type WeeklyParams struct {
	TSCode    string // 股票代码
	TradeDate string // 交易日期（YYYYMMDD）
	StartDate string // 开始日期(YYYYMMDD)
	EndDate   string // 结束日期(YYYYMMDD)
	Fields    string // 返回字段
}

// Weekly 获取周线行情（自动处理分页）
func (c *Client) Weekly(params *WeeklyParams, opts ...QueryOption) (*Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.TradeDate != "" {
		reqParams["trade_date"] = params.TradeDate
	}
	if params.StartDate != "" {
		reqParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		reqParams["end_date"] = params.EndDate
	}

	return c.Query("weekly", reqParams, params.Fields, opts...)
}

// ==================== 月线行情 ====================

// MonthlyParams 月线行情参数
type MonthlyParams struct {
	TSCode    string // 股票代码
	TradeDate string // 交易日期（YYYYMMDD）
	StartDate string // 开始日期(YYYYMMDD)
	EndDate   string // 结束日期(YYYYMMDD)
	Fields    string // 返回字段
}

// Monthly 获取月线行情（自动处理分页）
func (c *Client) Monthly(params *MonthlyParams, opts ...QueryOption) (*Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.TradeDate != "" {
		reqParams["trade_date"] = params.TradeDate
	}
	if params.StartDate != "" {
		reqParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		reqParams["end_date"] = params.EndDate
	}

	return c.Query("monthly", reqParams, params.Fields, opts...)
}

// ==================== 每日指标 ====================

// DailyBasicParams 每日指标参数
type DailyBasicParams struct {
	TSCode    string // 股票代码（二选一）
	TradeDate string // 交易日期（二选一）
	StartDate string // 开始日期
	EndDate   string // 结束日期
	Fields    string // 返回字段
}

// DailyBasic 获取每日指标（自动处理分页）
func (c *Client) DailyBasic(params *DailyBasicParams, opts ...QueryOption) (*Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.TradeDate != "" {
		reqParams["trade_date"] = params.TradeDate
	}
	if params.StartDate != "" {
		reqParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		reqParams["end_date"] = params.EndDate
	}

	return c.Query("daily_basic", reqParams, params.Fields, opts...)
}

// ==================== 个股资金流向 ====================

// MoneyFlowParams 个股资金流向参数
type MoneyFlowParams struct {
	TSCode    string // 股票代码
	TradeDate string // 交易日期
	StartDate string // 开始日期
	EndDate   string // 结束日期
	Fields    string // 返回字段
}

// MoneyFlow 获取个股资金流向（自动处理分页）
func (c *Client) MoneyFlow(params *MoneyFlowParams, opts ...QueryOption) (*Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.TradeDate != "" {
		reqParams["trade_date"] = params.TradeDate
	}
	if params.StartDate != "" {
		reqParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		reqParams["end_date"] = params.EndDate
	}

	return c.Query("moneyflow", reqParams, params.Fields, opts...)
}

// ==================== 财务数据接口 ====================

// IncomeParams 利润表参数
type IncomeParams struct {
	TSCode     string // 股票代码
	AnnDate    string // 公告日期（YYYYMMDD）
	StartDate  string // 公告开始日期
	EndDate    string // 公告结束日期
	Period     string // 报告期(每个季度最后一天的日期，比如20171231表示年报)
	ReportType string // 报告类型
	CompType   string // 公司类型
	Fields     string // 返回字段
}

// Income 获取利润表数据（自动处理分页）
func (c *Client) Income(params *IncomeParams, opts ...QueryOption) (*Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.AnnDate != "" {
		reqParams["ann_date"] = params.AnnDate
	}
	if params.StartDate != "" {
		reqParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		reqParams["end_date"] = params.EndDate
	}
	if params.Period != "" {
		reqParams["period"] = params.Period
	}
	if params.ReportType != "" {
		reqParams["report_type"] = params.ReportType
	}
	if params.CompType != "" {
		reqParams["comp_type"] = params.CompType
	}

	return c.Query("income", reqParams, params.Fields, opts...)
}

// BalanceSheetParams 资产负债表参数
type BalanceSheetParams struct {
	TSCode     string // 股票代码
	AnnDate    string // 公告日期
	StartDate  string // 公告开始日期
	EndDate    string // 公告结束日期
	Period     string // 报告期
	ReportType string // 报告类型
	CompType   string // 公司类型
	Fields     string // 返回字段
}

// BalanceSheet 获取资产负债表数据（自动处理分页）
func (c *Client) BalanceSheet(params *BalanceSheetParams, opts ...QueryOption) (*Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.AnnDate != "" {
		reqParams["ann_date"] = params.AnnDate
	}
	if params.StartDate != "" {
		reqParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		reqParams["end_date"] = params.EndDate
	}
	if params.Period != "" {
		reqParams["period"] = params.Period
	}
	if params.ReportType != "" {
		reqParams["report_type"] = params.ReportType
	}
	if params.CompType != "" {
		reqParams["comp_type"] = params.CompType
	}

	return c.Query("balancesheet", reqParams, params.Fields, opts...)
}

// CashFlowParams 现金流量表参数
type CashFlowParams struct {
	TSCode     string // 股票代码
	AnnDate    string // 公告日期
	StartDate  string // 公告开始日期
	EndDate    string // 公告结束日期
	Period     string // 报告期
	ReportType string // 报告类型
	CompType   string // 公司类型
	Fields     string // 返回字段
}

// CashFlow 获取现金流量表数据（自动处理分页）
func (c *Client) CashFlow(params *CashFlowParams, opts ...QueryOption) (*Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.AnnDate != "" {
		reqParams["ann_date"] = params.AnnDate
	}
	if params.StartDate != "" {
		reqParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		reqParams["end_date"] = params.EndDate
	}
	if params.Period != "" {
		reqParams["period"] = params.Period
	}
	if params.ReportType != "" {
		reqParams["report_type"] = params.ReportType
	}
	if params.CompType != "" {
		reqParams["comp_type"] = params.CompType
	}

	return c.Query("cashflow", reqParams, params.Fields, opts...)
}

// ==================== 指数数据 ====================

// IndexDailyParams 指数日线参数
type IndexDailyParams struct {
	TSCode    string // 指数代码
	TradeDate string // 交易日期
	StartDate string // 开始日期
	EndDate   string // 结束日期
	Fields    string // 返回字段
}

// IndexDaily 获取指数日线行情（自动处理分页）
func (c *Client) IndexDaily(params *IndexDailyParams, opts ...QueryOption) (*Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.TradeDate != "" {
		reqParams["trade_date"] = params.TradeDate
	}
	if params.StartDate != "" {
		reqParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		reqParams["end_date"] = params.EndDate
	}

	return c.Query("index_daily", reqParams, params.Fields, opts...)
}

// IndexBasicParams 指数基本信息参数
type IndexBasicParams struct {
	Market    string // 交易所或服务商代码
	Publisher string // 发布商
	Category  string // 指数类别
	Fields    string // 返回字段
}

// IndexBasic 获取指数基本信息（自动处理分页）
func (c *Client) IndexBasic(params *IndexBasicParams, opts ...QueryOption) (*Response, error) {
	reqParams := make(map[string]interface{})
	if params.Market != "" {
		reqParams["market"] = params.Market
	}
	if params.Publisher != "" {
		reqParams["publisher"] = params.Publisher
	}
	if params.Category != "" {
		reqParams["category"] = params.Category
	}

	return c.Query("index_basic", reqParams, params.Fields, opts...)
}

// ==================== 期货数据 ====================

// FutBasicParams 期货合约信息参数
type FutBasicParams struct {
	Exchange string // 交易所代码
	FutType  string // 合约类型
	Fields   string // 返回字段
}

// FutBasic 获取期货合约信息（自动处理分页）
func (c *Client) FutBasic(params *FutBasicParams, opts ...QueryOption) (*Response, error) {
	reqParams := make(map[string]interface{})
	if params.Exchange != "" {
		reqParams["exchange"] = params.Exchange
	}
	if params.FutType != "" {
		reqParams["fut_type"] = params.FutType
	}

	return c.Query("fut_basic", reqParams, params.Fields, opts...)
}

// FutDailyParams 期货日线参数
type FutDailyParams struct {
	TSCode    string // 期货合约代码
	TradeDate string // 交易日期
	StartDate string // 开始日期
	EndDate   string // 结束日期
	Exchange  string // 交易所
	Fields    string // 返回字段
}

// FutDaily 获取期货日线行情（自动处理分页）
func (c *Client) FutDaily(params *FutDailyParams, opts ...QueryOption) (*Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.TradeDate != "" {
		reqParams["trade_date"] = params.TradeDate
	}
	if params.StartDate != "" {
		reqParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		reqParams["end_date"] = params.EndDate
	}
	if params.Exchange != "" {
		reqParams["exchange"] = params.Exchange
	}

	return c.Query("fut_daily", reqParams, params.Fields, opts...)
}

// ==================== 辅助工具方法 ====================

// String 辅助函数：将 interface{} 转换为字符串
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}
