// Package basic 提供 Tushare 股票基础数据相关接口
// 文档参考: https://tushare.pro/document/2?doc_id=25
package basic

import (
	tushare "github.com/fletcherlau/go-tushare"
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
