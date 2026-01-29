// Package basic 提供 Tushare 股票基础数据相关接口
// 文档参考: https://tushare.pro/document/2?doc_id=25
package basic

import (
	"strings"

	tushare "github.com/fletcherlau/go-tushare"
)

// Exchange 交易所代码
type Exchange string

const (
	// ExchangeSSE 上交所
	ExchangeSSE Exchange = "SSE"
	// ExchangeSZSE 深交所
	ExchangeSZSE Exchange = "SZSE"
	// ExchangeBSE 北交所
	ExchangeBSE Exchange = "BSE"
)

// Market 市场类别
type Market string

const (
	// MarketMain 主板
	MarketMain Market = "主板"
	// MarketGEM 创业板
	MarketGEM Market = "创业板"
	// MarketSTAR 科创板
	MarketSTAR Market = "科创板"
	// MarketCDR CDR
	MarketCDR Market = "CDR"
	// MarketBSE 北交所
	MarketBSE Market = "北交所"
)

// IsHS 是否沪深港通标的
type IsHS string

const (
	// IsHSNo 否
	IsHSNo IsHS = "N"
	// IsHSShanghai 沪股通
	IsHSShanghai IsHS = "H"
	// IsHSShenzhen 深股通
	IsHSShenzhen IsHS = "S"
)

// ListStatus 上市状态
type ListStatus string

const (
	// ListStatusListed 上市
	ListStatusListed ListStatus = "L"
	// ListStatusDelisted 退市
	ListStatusDelisted ListStatus = "D"
	// ListStatusSuspended 暂停上市
	ListStatusSuspended ListStatus = "P"
)

// StockBasicField 返回字段常量
const (
	StockBasicFieldTSCode     = "ts_code"     // TS代码
	StockBasicFieldSymbol     = "symbol"      // 股票代码
	StockBasicFieldName       = "name"        // 股票名称
	StockBasicFieldArea       = "area"        // 地域
	StockBasicFieldIndustry   = "industry"    // 所属行业
	StockBasicFieldFullName   = "fullname"    // 股票全称
	StockBasicFieldEnName     = "enname"      // 英文全称
	StockBasicFieldCNSpell    = "cnspell"     // 拼音缩写
	StockBasicFieldMarket     = "market"      // 市场类型
	StockBasicFieldExchange   = "exchange"    // 交易所代码
	StockBasicFieldCurrType   = "curr_type"   // 交易货币
	StockBasicFieldListStatus = "list_status" // 上市状态
	StockBasicFieldListDate   = "list_date"   // 上市日期
	StockBasicFieldDelistDate = "delist_date" // 退市日期
	StockBasicFieldIsHS       = "is_hs"       // 是否沪深港通标的
)

// StockBasicParams 股票基础信息参数
// 接口: stock_basic
// 描述: 获取股票基础信息，包括股票代码、名称、上市日期、退市日期等
type StockBasicParams struct {
	TSCode     string     // TS股票代码，支持单个或多个（逗号分隔）
	Name       string     // 股票名称
	Exchange   Exchange   // 交易所代码
	Market     Market     // 市场类别
	IsHS       IsHS       // 是否沪深港通标的
	ListStatus ListStatus // 上市状态，默认L
	Fields     []string   // 返回字段列表
}

// StockBasicItem 股票基础信息响应项
type StockBasicItem struct {
	TSCode     string     `json:"ts_code"`     // TS代码
	Symbol     string     `json:"symbol"`      // 股票代码
	Name       string     `json:"name"`        // 股票名称
	Area       string     `json:"area"`        // 地域
	Industry   string     `json:"industry"`    // 所属行业
	FullName   string     `json:"fullname"`    // 股票全称
	EnName     string     `json:"enname"`      // 英文全称
	CNSpell    string     `json:"cnspell"`     // 拼音缩写
	Market     Market     `json:"market"`      // 市场类型
	Exchange   Exchange   `json:"exchange"`    // 交易所代码
	CurrType   string     `json:"curr_type"`   // 交易货币
	ListStatus ListStatus `json:"list_status"` // 上市状态
	ListDate   string     `json:"list_date"`   // 上市日期
	DelistDate string     `json:"delist_date"` // 退市日期
	IsHS       IsHS       `json:"is_hs"`       // 是否沪深港通标的
}

// StockBasic 获取股票基础信息（自动处理分页）
// 根据指定条件获取股票基础信息数据
func StockBasic(c *tushare.Client, params *StockBasicParams, opts ...tushare.QueryOption) ([]*StockBasicItem, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.Name != "" {
		reqParams["name"] = params.Name
	}
	if params.Exchange != "" {
		reqParams["exchange"] = string(params.Exchange)
	}
	if params.Market != "" {
		reqParams["market"] = string(params.Market)
	}
	if params.IsHS != "" {
		reqParams["is_hs"] = string(params.IsHS)
	}
	if params.ListStatus != "" {
		reqParams["list_status"] = string(params.ListStatus)
	} else {
		reqParams["list_status"] = string(ListStatusListed)
	}

	fields := ""
	if len(params.Fields) > 0 {
		fields = strings.Join(params.Fields, ",")
	}

	resp, err := c.Query("stock_basic", reqParams, fields, opts...)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, &tushare.APIError{
			Code: resp.Code,
			Msg:  resp.Msg,
		}
	}

	var items []*StockBasicItem
	if err := resp.ToStruct(&items); err != nil {
		return nil, err
	}

	return items, nil
}
