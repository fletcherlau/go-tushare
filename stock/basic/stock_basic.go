// Package basic 提供 Tushare 股票基础数据相关接口
// 文档参考: https://tushare.pro/document/2?doc_id=25
package basic

import (
	tushare "github.com/yourusername/go-tushare"
)

// StockBasicParams 股票基础信息参数
// 接口: stock_basic
// 描述: 获取股票基础信息，包括股票代码、名称、上市日期、退市日期等
type StockBasicParams struct {
	TSCode     string // TS股票代码，支持单个或多个（逗号分隔）
	Name       string // 股票名称
	Exchange   string // 交易所代码：SSE上交所 SZSE深交所 BSE北交所
	Market     string // 市场类别：主板/创业板/科创板/CDR/北交所
	IsHS       string // 是否沪深港通标的：N否 H沪股通 S深股通
	ListStatus string // 上市状态：L上市 D退市 P暂停上市，默认L
	Fields     string // 返回字段，用逗号分隔
}

// StockBasic 获取股票基础信息（自动处理分页）
// 根据指定条件获取股票基础信息数据
func StockBasic(c *tushare.Client, params *StockBasicParams, opts ...tushare.QueryOption) (*tushare.Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.Name != "" {
		reqParams["name"] = params.Name
	}
	if params.Exchange != "" {
		reqParams["exchange"] = params.Exchange
	}
	if params.Market != "" {
		reqParams["market"] = params.Market
	}
	if params.IsHS != "" {
		reqParams["is_hs"] = params.IsHS
	}
	if params.ListStatus != "" {
		reqParams["list_status"] = params.ListStatus
	} else {
		reqParams["list_status"] = "L"
	}

	return c.Query("stock_basic", reqParams, params.Fields, opts...)
}

// StockCompanyParams 上市公司基本信息参数
// 接口: stock_company
// 描述: 获取上市公司基本信息，包括公司全称、注册地、主营业务等
type StockCompanyParams struct {
	TSCode   string // TS股票代码
	Exchange string // 交易所代码：SZSE深交所 SSE上交所
	Fields   string // 返回字段，用逗号分隔
}

// StockCompany 获取上市公司基本信息（自动处理分页）
func StockCompany(c *tushare.Client, params *StockCompanyParams, opts ...tushare.QueryOption) (*tushare.Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.Exchange != "" {
		reqParams["exchange"] = params.Exchange
	}

	return c.Query("stock_company", reqParams, params.Fields, opts...)
}

// TradeCalParams 交易日历参数
// 接口: trade_cal
// 描述: 获取各大交易所交易日历数据
type TradeCalParams struct {
	Exchange  string // 交易所代码：SSE上交所 SZSE深交所
	StartDate string // 开始日期（YYYYMMDD）
	EndDate   string // 结束日期（YYYYMMDD）
	IsOpen    string // 是否交易：0休市 1交易
	Fields    string // 返回字段
}

// TradeCal 获取交易日历（自动处理分页）
func TradeCal(c *tushare.Client, params *TradeCalParams, opts ...tushare.QueryOption) (*tushare.Response, error) {
	reqParams := make(map[string]interface{})
	if params.Exchange != "" {
		reqParams["exchange"] = params.Exchange
	}
	if params.StartDate != "" {
		reqParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		reqParams["end_date"] = params.EndDate
	}
	if params.IsOpen != "" {
		reqParams["is_open"] = params.IsOpen
	}
	return c.Query("trade_cal", reqParams, params.Fields, opts...)
}

// NameChangeParams 股票曾用名参数
// 接口: namechange
// 描述: 获取股票曾用名信息
type NameChangeParams struct {
	TSCode    string // TS股票代码
	StartDate string // 公告开始日期（YYYYMMDD）
	EndDate   string // 公告结束日期（YYYYMMDD）
	Fields    string // 返回字段
}

// NameChange 获取股票曾用名（自动处理分页）
func NameChange(c *tushare.Client, params *NameChangeParams, opts ...tushare.QueryOption) (*tushare.Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.StartDate != "" {
		reqParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		reqParams["end_date"] = params.EndDate
	}
	return c.Query("namechange", reqParams, params.Fields, opts...)
}

// HSConstParams 沪深股通成份股参数
// 接口: hs_const
// 描述: 获取沪深港通成份股信息
type HSConstParams struct {
	HsType string // 类型：SH沪股通 SZ深股通
	IsNew  string // 是否最新：1是 0否（默认1）
	Fields string // 返回字段
}

// HSConst 获取沪深股通成份股（自动处理分页）
func HSConst(c *tushare.Client, params *HSConstParams, opts ...tushare.QueryOption) (*tushare.Response, error) {
	reqParams := make(map[string]interface{})
	if params.HsType != "" {
		reqParams["hs_type"] = params.HsType
	}
	if params.IsNew != "" {
		reqParams["is_new"] = params.IsNew
	}
	return c.Query("hs_const", reqParams, params.Fields, opts...)
}

// StockSuspendParams 停牌信息参数
// 接口: suspend
// 描述: 获取股票停牌信息
type StockSuspendParams struct {
	TSCode       string // TS股票代码
	SuspendDate  string // 停牌日期（YYYYMMDD）
	ResumeDate   string // 复牌日期（YYYYMMDD）
	Fields       string // 返回字段
}

// StockSuspend 获取停牌信息（自动处理分页）
func StockSuspend(c *tushare.Client, params *StockSuspendParams, opts ...tushare.QueryOption) (*tushare.Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.SuspendDate != "" {
		reqParams["suspend_date"] = params.SuspendDate
	}
	if params.ResumeDate != "" {
		reqParams["resume_date"] = params.ResumeDate
	}
	return c.Query("suspend", reqParams, params.Fields, opts...)
}

// StkLimitParams 个股涨跌停参数
// 接口: stk_limit
// 描述: 获取个股涨跌停价格及状态
type StkLimitParams struct {
	TSCode    string // TS股票代码
	TradeDate string // 交易日期（YYYYMMDD）
	StartDate string // 开始日期（YYYYMMDD）
	EndDate   string // 结束日期（YYYYMMDD）
	Fields    string // 返回字段
}

// StkLimit 获取个股涨跌停（自动处理分页）
func StkLimit(c *tushare.Client, params *StkLimitParams, opts ...tushare.QueryOption) (*tushare.Response, error) {
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
	return c.Query("stk_limit", reqParams, params.Fields, opts...)
}

// StkRewardParams 股票质押参数
// 接口: stk_reward
// 描述: 获取上市公司员工持股计划和股权激励信息
type StkRewardParams struct {
	TSCode string // TS股票代码
	AnnDate string // 公告日期（YYYYMMDD）
	Fields string // 返回字段
}

// StkReward 获取股票质押信息（自动处理分页）
func StkReward(c *tushare.Client, params *StkRewardParams, opts ...tushare.QueryOption) (*tushare.Response, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.AnnDate != "" {
		reqParams["ann_date"] = params.AnnDate
	}
	return c.Query("stk_reward", reqParams, params.Fields, opts...)
}
