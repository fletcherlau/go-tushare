package financial

import (
	"strings"

	tushare "github.com/fletcherlau/go-tushare"
)

// CompType 公司类型
type CompType string

const (
	// CompTypeGeneral 一般工商业
	CompTypeGeneral CompType = "1"
	// CompTypeBank 银行
	CompTypeBank CompType = "2"
	// CompTypeInsurance 保险
	CompTypeInsurance CompType = "3"
	// CompTypeSecurities 证券
	CompTypeSecurities CompType = "4"
)

// IncomeField 返回字段常量
const (
	IncomeFieldTSCode       = "ts_code"        // TS代码
	IncomeFieldAnnDate      = "ann_date"       // 公告日期
	IncomeFieldFAnnDate     = "f_ann_date"     // 实际公告日期
	IncomeFieldEndDate      = "end_date"       // 报告期
	IncomeFieldReportType   = "report_type"    // 报告类型
	IncomeFieldCompType     = "comp_type"      // 公司类型
	IncomeFieldBasicEps     = "basic_eps"      // 基本每股收益
	IncomeFieldDilutedEps   = "diluted_eps"    // 稀释每股收益
	IncomeFieldTotalRevenue = "total_revenue"  // 营业总收入
	IncomeFieldRevenue      = "revenue"        // 营业收入
	IncomeFieldOperateProfit = "operate_profit" // 营业利润
	IncomeFieldTotalProfit  = "total_profit"   // 利润总额
	IncomeFieldNIncome      = "n_income"       // 净利润（含少数股东损益）
	IncomeFieldEBIT         = "ebit"           // 息税前利润
	IncomeFieldEBITDA       = "ebitda"         // 息税折旧摊销前利润
)

// IncomeParams 利润表参数
// 接口: income
// 描述: 获取上市公司财务利润表数据
// 文档: https://tushare.pro/document/2?doc_id=33
type IncomeParams struct {
	TSCode      string   // 股票代码（必填）
	AnnDate     string   // 公告日期（YYYYMMDD格式）
	FAnnDate    string   // 实际公告日期
	StartDate   string   // 公告日开始日期
	EndDate     string   // 公告日结束日期
	Period      string   // 报告期（每个季度最后一天的日期）
	ReportType  string   // 报告类型
	CompType    CompType // 公司类型（1一般工商业 2银行 3保险 4证券）
	Fields      []string // 返回字段列表
}

// IncomeItem 利润表响应项
type IncomeItem struct {
	TSCode       string  `json:"ts_code"`        // TS代码
	AnnDate      string  `json:"ann_date"`       // 公告日期
	FAnnDate     string  `json:"f_ann_date"`     // 实际公告日期
	EndDate      string  `json:"end_date"`       // 报告期
	ReportType   string  `json:"report_type"`    // 报告类型
	CompType     string  `json:"comp_type"`      // 公司类型
	BasicEps     float64 `json:"basic_eps"`      // 基本每股收益
	DilutedEps   float64 `json:"diluted_eps"`    // 稀释每股收益
	TotalRevenue float64 `json:"total_revenue"`  // 营业总收入
	Revenue      float64 `json:"revenue"`        // 营业收入
	OperateProfit float64 `json:"operate_profit"` // 营业利润
	TotalProfit  float64 `json:"total_profit"`   // 利润总额
	NIncome      float64 `json:"n_income"`       // 净利润（含少数股东损益）
	EBIT         float64 `json:"ebit"`           // 息税前利润
	EBITDA       float64 `json:"ebitda"`         // 息税折旧摊销前利润
}

// Income 获取利润表数据（自动处理分页）
func Income(c *tushare.Client, params *IncomeParams, opts ...tushare.QueryOption) ([]*IncomeItem, error) {
	reqParams := make(map[string]interface{})
	if params.TSCode != "" {
		reqParams["ts_code"] = params.TSCode
	}
	if params.AnnDate != "" {
		reqParams["ann_date"] = params.AnnDate
	}
	if params.FAnnDate != "" {
		reqParams["f_ann_date"] = params.FAnnDate
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
		reqParams["comp_type"] = string(params.CompType)
	}

	fields := ""
	if len(params.Fields) > 0 {
		fields = strings.Join(params.Fields, ",")
	}

	resp, err := c.Query("income", reqParams, fields, opts...)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, &tushare.APIError{
			Code: resp.Code,
			Msg:  resp.Msg,
		}
	}

	var items []*IncomeItem
	if err := resp.ToStruct(&items); err != nil {
		return nil, err
	}

	return items, nil
}
