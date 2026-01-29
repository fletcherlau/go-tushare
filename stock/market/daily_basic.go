package market

import (
	"strings"

	tushare "github.com/fletcherlau/go-tushare"
)

// DailyBasicField 返回字段常量
const (
	DailyBasicFieldTSCode        = "ts_code"         // TS股票代码
	DailyBasicFieldTradeDate     = "trade_date"      // 交易日期
	DailyBasicFieldClose         = "close"           // 当日收盘价
	DailyBasicFieldTurnoverRate  = "turnover_rate"   // 换手率（%）
	DailyBasicFieldTurnoverRateF = "turnover_rate_f" // 换手率（自由流通股）
	DailyBasicFieldVolumeRatio   = "volume_ratio"    // 量比
	DailyBasicFieldPE            = "pe"              // 市盈率（总市值/净利润，亏损的PE为空）
	DailyBasicFieldPETTM         = "pe_ttm"          // 市盈率（TTM，亏损的PE为空）
	DailyBasicFieldPB            = "pb"              // 市净率（总市值/净资产）
	DailyBasicFieldPS            = "ps"              // 市销率
	DailyBasicFieldPSTTM         = "ps_ttm"          // 市销率（TTM）
	DailyBasicFieldDVRatio       = "dv_ratio"        // 股息率（%）
	DailyBasicFieldDVTTM         = "dv_ttm"          // 股息率（TTM）（%）
	DailyBasicFieldTotalShare    = "total_share"     // 总股本（万股）
	DailyBasicFieldFloatShare    = "float_share"     // 流通股本（万股）
	DailyBasicFieldFreeShare     = "free_share"      // 自由流通股本（万）
	DailyBasicFieldTotalMV       = "total_mv"        // 总市值（万元）
	DailyBasicFieldCircMV        = "circ_mv"         // 流通市值（万元）
)

// DailyBasicParams 每日指标参数
// 接口: daily_basic
// 描述: 获取全部股票每日重要的基本面指标，可用于选股分析、报表展示等。单次请求最大返回6000条数据，可按日线循环提取全部历史。
// 文档: https://tushare.pro/document/2?doc_id=32
type DailyBasicParams struct {
	TSCode     string   // 股票代码（二选一）
	TradeDate  string   // 交易日期（二选一）
	StartDate  string   // 开始日期(YYYYMMDD)
	EndDate    string   // 结束日期(YYYYMMDD)
	Fields     []string // 返回字段列表
}

// DailyBasicItem 每日指标响应项
type DailyBasicItem struct {
	TSCode        string  `json:"ts_code"`         // TS股票代码
	TradeDate     string  `json:"trade_date"`      // 交易日期
	Close         float64 `json:"close"`           // 当日收盘价
	TurnoverRate  float64 `json:"turnover_rate"`   // 换手率（%）
	TurnoverRateF float64 `json:"turnover_rate_f"` // 换手率（自由流通股）
	VolumeRatio   float64 `json:"volume_ratio"`    // 量比
	PE            float64 `json:"pe"`              // 市盈率（总市值/净利润，亏损的PE为空）
	PETTM         float64 `json:"pe_ttm"`          // 市盈率（TTM，亏损的PE为空）
	PB            float64 `json:"pb"`              // 市净率（总市值/净资产）
	PS            float64 `json:"ps"`              // 市销率
	PSTTM         float64 `json:"ps_ttm"`          // 市销率（TTM）
	DVRatio       float64 `json:"dv_ratio"`        // 股息率（%）
	DVTTM         float64 `json:"dv_ttm"`          // 股息率（TTM）（%）
	TotalShare    float64 `json:"total_share"`     // 总股本（万股）
	FloatShare    float64 `json:"float_share"`     // 流通股本（万股）
	FreeShare     float64 `json:"free_share"`      // 自由流通股本（万）
	TotalMV       float64 `json:"total_mv"`        // 总市值（万元）
	CircMV        float64 `json:"circ_mv"`         // 流通市值（万元）
}

// DailyBasic 获取每日指标数据（自动处理分页）
func DailyBasic(c *tushare.Client, params *DailyBasicParams, opts ...tushare.QueryOption) ([]*DailyBasicItem, error) {
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

	resp, err := c.Query("daily_basic", reqParams, fields, opts...)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, &tushare.APIError{
			Code: resp.Code,
			Msg:  resp.Msg,
		}
	}

	var items []*DailyBasicItem
	if err := resp.ToStruct(&items); err != nil {
		return nil, err
	}

	return items, nil
}
