// Package market 提供 Tushare 股票行情数据相关接口
// 文档参考: https://tushare.pro/document/2?doc_id=27
package market

import (
	"strings"

	tushare "github.com/fletcherlau/go-tushare"
)

// DailyField 返回字段常量
const (
	DailyFieldTSCode     = "ts_code"     // 股票代码
	DailyFieldTradeDate  = "trade_date"  // 交易日期
	DailyFieldOpen       = "open"        // 开盘价
	DailyFieldHigh       = "high"        // 最高价
	DailyFieldLow        = "low"         // 最低价
	DailyFieldClose      = "close"       // 收盘价
	DailyFieldPreClose   = "pre_close"   // 昨收价【除权价】
	DailyFieldChange     = "change"      // 涨跌额
	DailyFieldPctChg     = "pct_chg"     // 涨跌幅
	DailyFieldVol        = "vol"         // 成交量（手）
	DailyFieldAmount     = "amount"      // 成交额（千元）
)

// DailyParams A股日线行情参数
// 接口: daily
// 描述: 获取股票行情数据。交易日每天15点～16点之间入库。本接口是未复权行情，停牌期间不提供数据。
// 调用限制：基础积分每分钟内可调取500次，每次6000条数据。
// 文档: https://tushare.pro/document/2?doc_id=27
type DailyParams struct {
	TSCode     string   // 股票代码（支持多个股票同时提取，逗号分隔）
	TradeDate  string   // 交易日期（YYYYMMDD）
	StartDate  string   // 开始日期(YYYYMMDD)
	EndDate    string   // 结束日期(YYYYMMDD)
	Fields     []string // 返回字段列表
}

// DailyItem A股日线行情响应项
type DailyItem struct {
	TSCode     string  `json:"ts_code"`     // 股票代码
	TradeDate  string  `json:"trade_date"`  // 交易日期
	Open       float64 `json:"open"`        // 开盘价
	High       float64 `json:"high"`        // 最高价
	Low        float64 `json:"low"`         // 最低价
	Close      float64 `json:"close"`       // 收盘价
	PreClose   float64 `json:"pre_close"`   // 昨收价【除权价】
	Change     float64 `json:"change"`      // 涨跌额
	PctChg     float64 `json:"pct_chg"`     // 涨跌幅
	Vol        float64 `json:"vol"`         // 成交量（手）
	Amount     float64 `json:"amount"`      // 成交额（千元）
}

// Daily 获取A股日线行情数据（自动处理分页）
// 根据指定条件获取股票的日线行情数据
func Daily(c *tushare.Client, params *DailyParams, opts ...tushare.QueryOption) ([]*DailyItem, error) {
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

	fields := ""
	if len(params.Fields) > 0 {
		fields = strings.Join(params.Fields, ",")
	}

	resp, err := c.Query("daily", reqParams, fields, opts...)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, &tushare.APIError{
			Code: resp.Code,
			Msg:  resp.Msg,
		}
	}

	var items []*DailyItem
	if err := resp.ToStruct(&items); err != nil {
		return nil, err
	}

	return items, nil
}
