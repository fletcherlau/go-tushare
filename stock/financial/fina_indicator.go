package financial

import (
	"strings"

	tushare "github.com/fletcherlau/go-tushare"
)

// FinaIndicatorField 返回字段常量（核心财务指标）
const (
	FinaIndicatorFieldTSCode          = "ts_code"           // TS代码
	FinaIndicatorFieldAnnDate         = "ann_date"          // 公告日期
	FinaIndicatorFieldEndDate         = "end_date"          // 报告期
	FinaIndicatorFieldEps             = "eps"               // 基本每股收益
	FinaIndicatorFieldDtEps           = "dt_eps"            // 稀释每股收益
	FinaIndicatorFieldTotalRevenuePs  = "total_revenue_ps"  // 每股营业总收入
	FinaIndicatorFieldRevenuePs       = "revenue_ps"        // 每股营业收入
	FinaIndicatorFieldCapitalResePs   = "capital_rese_ps"   // 每股资本公积
	FinaIndicatorFieldSurplusResePs   = "surplus_rese_ps"   // 每股盈余公积
	FinaIndicatorFieldUndistProfitPs  = "undist_profit_ps"  // 每股未分配利润
	FinaIndicatorFieldExtraItem       = "extra_item"        // 非经常性损益
	FinaIndicatorFieldProfitDedt      = "profit_dedt"       // 扣除非经常性损益后的净利润（扣非净利润）
	FinaIndicatorFieldGrossMargin     = "gross_margin"      // 毛利
	FinaIndicatorFieldCurrentRatio    = "current_ratio"     // 流动比率
	FinaIndicatorFieldQuickRatio      = "quick_ratio"       // 速动比率
	FinaIndicatorFieldCashRatio       = "cash_ratio"        // 保守速动比率
	FinaIndicatorFieldRoe             = "roe"               // 净资产收益率
	FinaIndicatorFieldRoa             = "roa"               // 总资产报酬率
	FinaIndicatorFieldDebtToAssets    = "debt_to_assets"    // 资产负债率
	FinaIndicatorFieldBasicEpsYoy     = "basic_eps_yoy"     // 基本每股收益同比增长率（%）
	FinaIndicatorFieldNetprofitYoy    = "netprofit_yoy"     // 归属母公司股东的净利润同比增长率（%）
	FinaIndicatorFieldOcfYoy          = "ocf_yoy"           // 经营活动产生的现金流量净额同比增长率（%）
)

// FinaIndicatorParams 财务指标参数
// 接口: fina_indicator
// 描述: 获取上市公司财务指标数据
// 注意: 该接口返回字段较多（100+个），为避免服务器压力，每次请求最多返回100条记录
// 文档: https://tushare.pro/document/2?doc_id=79
type FinaIndicatorParams struct {
	TSCode     string   // TS股票代码，如 600001.SH/000001.SZ（必填）
	AnnDate    string   // 公告日期
	StartDate  string   // 报告期开始日期
	EndDate    string   // 报告期结束日期
	Period     string   // 报告期（每个季度最后一天的日期，如20171231表示年报）
	Fields     []string // 返回字段列表
}

// FinaIndicatorItem 财务指标响应项（核心字段）
// 注意: 实际接口返回100+个字段，这里包含核心财务指标
type FinaIndicatorItem struct {
	TSCode          string  `json:"ts_code"`           // TS代码
	AnnDate         string  `json:"ann_date"`          // 公告日期
	EndDate         string  `json:"end_date"`          // 报告期
	Eps             float64 `json:"eps"`               // 基本每股收益
	DtEps           float64 `json:"dt_eps"`            // 稀释每股收益
	TotalRevenuePs  float64 `json:"total_revenue_ps"`  // 每股营业总收入
	RevenuePs       float64 `json:"revenue_ps"`        // 每股营业收入
	CapitalResePs   float64 `json:"capital_rese_ps"`   // 每股资本公积
	SurplusResePs   float64 `json:"surplus_rese_ps"`   // 每股盈余公积
	UndistProfitPs  float64 `json:"undist_profit_ps"`  // 每股未分配利润
	ExtraItem       float64 `json:"extra_item"`        // 非经常性损益
	ProfitDedt      float64 `json:"profit_dedt"`       // 扣除非经常性损益后的净利润（扣非净利润）
	GrossMargin     float64 `json:"gross_margin"`      // 毛利
	CurrentRatio    float64 `json:"current_ratio"`     // 流动比率
	QuickRatio      float64 `json:"quick_ratio"`       // 速动比率
	CashRatio       float64 `json:"cash_ratio"`        // 保守速动比率
	Roe             float64 `json:"roe"`               // 净资产收益率
	Roa             float64 `json:"roa"`               // 总资产报酬率
	DebtToAssets    float64 `json:"debt_to_assets"`    // 资产负债率
	BasicEpsYoy     float64 `json:"basic_eps_yoy"`     // 基本每股收益同比增长率（%）
	NetprofitYoy    float64 `json:"netprofit_yoy"`     // 归属母公司股东的净利润同比增长率（%）
	OcfYoy          float64 `json:"ocf_yoy"`           // 经营活动产生的现金流量净额同比增长率（%）
}

// FinaIndicator 获取财务指标数据（自动处理分页）
// 注意: 该接口每次请求最多返回100条记录
func FinaIndicator(c *tushare.Client, params *FinaIndicatorParams, opts ...tushare.QueryOption) ([]*FinaIndicatorItem, error) {
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

	fields := ""
	if len(params.Fields) > 0 {
		fields = strings.Join(params.Fields, ",")
	}

	resp, err := c.Query("fina_indicator", reqParams, fields, opts...)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, &tushare.APIError{
			Code: resp.Code,
			Msg:  resp.Msg,
		}
	}

	var items []*FinaIndicatorItem
	if err := resp.ToStruct(&items); err != nil {
		return nil, err
	}

	return items, nil
}
