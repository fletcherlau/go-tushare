// Package basic 提供 Tushare 股票基础数据相关接口
// 文档参考: https://tushare.pro/document/2?doc_id=26
package basic

import (
	"strings"

	tushare "github.com/fletcherlau/go-tushare"
)

// TradeCalExchange 交易日历交易所代码
type TradeCalExchange string

const (
	// TradeCalExchangeSSE 上交所
	TradeCalExchangeSSE TradeCalExchange = "SSE"
	// TradeCalExchangeSZSE 深交所
	TradeCalExchangeSZSE TradeCalExchange = "SZSE"
	// TradeCalExchangeCFFEX 中金所
	TradeCalExchangeCFFEX TradeCalExchange = "CFFEX"
	// TradeCalExchangeSHFE 上期所
	TradeCalExchangeSHFE TradeCalExchange = "SHFE"
	// TradeCalExchangeCZCE 郑商所
	TradeCalExchangeCZCE TradeCalExchange = "CZCE"
	// TradeCalExchangeDCE 大商所
	TradeCalExchangeDCE TradeCalExchange = "DCE"
	// TradeCalExchangeINE 上能源
	TradeCalExchangeINE TradeCalExchange = "INE"
)

// TradeCalIsOpen 是否交易
type TradeCalIsOpen string

const (
	// TradeCalIsOpenNo 休市
	TradeCalIsOpenNo TradeCalIsOpen = "0"
	// TradeCalIsOpenYes 交易
	TradeCalIsOpenYes TradeCalIsOpen = "1"
)

// TradeCalField 返回字段常量
const (
	TradeCalFieldExchange     = "exchange"      // 交易所代码
	TradeCalFieldCalDate      = "cal_date"      // 日历日期
	TradeCalFieldIsOpen       = "is_open"       // 是否交易
	TradeCalFieldPretradeDate = "pretrade_date" // 上一个交易日
)

// TradeCalParams 交易日历参数
// 接口: trade_cal
// 描述: 获取各大交易所交易日历数据，默认提取的是上交所
// 文档: https://tushare.pro/document/2?doc_id=26
type TradeCalParams struct {
	Exchange  TradeCalExchange // 交易所代码，默认SSE
	StartDate string           // 开始日期（格式：YYYYMMDD）
	EndDate   string           // 结束日期（格式：YYYYMMDD）
	IsOpen    TradeCalIsOpen   // 是否交易：'0'表示休市，'1'表示交易
	Fields    []string         // 返回字段列表
}

// TradeCalItem 交易日历响应项
type TradeCalItem struct {
	Exchange     TradeCalExchange `json:"exchange"`      // 交易所代码
	CalDate      string           `json:"cal_date"`      // 日历日期
	IsOpen       TradeCalIsOpen   `json:"is_open"`       // 是否交易
	PretradeDate string           `json:"pretrade_date"` // 上一个交易日
}

// TradeCal 获取交易日历数据（自动处理分页）
// 根据指定条件获取各大交易所的交易日历信息
func TradeCal(c *tushare.Client, params *TradeCalParams, opts ...tushare.QueryOption) ([]*TradeCalItem, error) {
	reqParams := make(map[string]interface{})
	if params.Exchange != "" {
		reqParams["exchange"] = string(params.Exchange)
	} else {
		reqParams["exchange"] = string(TradeCalExchangeSSE)
	}
	if params.StartDate != "" {
		reqParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		reqParams["end_date"] = params.EndDate
	}
	if params.IsOpen != "" {
		reqParams["is_open"] = string(params.IsOpen)
	}

	fields := ""
	if len(params.Fields) > 0 {
		fields = strings.Join(params.Fields, ",")
	}

	resp, err := c.Query("trade_cal", reqParams, fields, opts...)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, &tushare.APIError{
			Code: resp.Code,
			Msg:  resp.Msg,
		}
	}

	var items []*TradeCalItem
	if err := resp.ToStruct(&items); err != nil {
		return nil, err
	}

	return items, nil
}
