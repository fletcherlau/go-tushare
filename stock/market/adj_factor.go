package market

import (
	"strings"

	tushare "github.com/fletcherlau/go-tushare"
)

// AdjFactorField 返回字段常量
const (
	AdjFactorFieldTSCode     = "ts_code"     // 股票代码
	AdjFactorFieldTradeDate  = "trade_date"  // 交易日期
	AdjFactorFieldAdjFactor  = "adj_factor"  // 复权因子
)

// AdjFactorParams 复权因子参数
// 接口: adj_factor
// 描述: 获取股票复权因子，可提取单只股票全部历史复权因子，也可以提取单日全部股票的复权因子
// 文档: https://tushare.pro/document/2?doc_id=28
type AdjFactorParams struct {
	TSCode     string   // 股票代码
	TradeDate  string   // 交易日期(YYYYMMDD)
	StartDate  string   // 开始日期
	EndDate    string   // 结束日期
	Fields     []string // 返回字段列表
}

// AdjFactorItem 复权因子响应项
type AdjFactorItem struct {
	TSCode     string  `json:"ts_code"`     // 股票代码
	TradeDate  string  `json:"trade_date"`  // 交易日期
	AdjFactor  float64 `json:"adj_factor"`  // 复权因子
}

// AdjFactor 获取复权因子数据（自动处理分页）
func AdjFactor(c *tushare.Client, params *AdjFactorParams, opts ...tushare.QueryOption) ([]*AdjFactorItem, error) {
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

	resp, err := c.Query("adj_factor", reqParams, fields, opts...)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, &tushare.APIError{
			Code: resp.Code,
			Msg:  resp.Msg,
		}
	}

	var items []*AdjFactorItem
	if err := resp.ToStruct(&items); err != nil {
		return nil, err
	}

	return items, nil
}
