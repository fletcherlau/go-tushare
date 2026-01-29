package financial

import (
	"strings"

	tushare "github.com/fletcherlau/go-tushare"
)

// CashFlowField 返回字段常量
const (
	CashFlowFieldTSCode               = "ts_code"                // TS股票代码
	CashFlowFieldAnnDate              = "ann_date"               // 公告日期
	CashFlowFieldFAnnDate             = "f_ann_date"             // 实际公告日期
	CashFlowFieldEndDate              = "end_date"               // 报告期
	CashFlowFieldCompType             = "comp_type"              // 公司类型
	CashFlowFieldReportType           = "report_type"            // 报表类型
	CashFlowFieldNetProfit            = "net_profit"             // 净利润
	CashFlowFieldFinanExp             = "finan_exp"              // 财务费用
	CashFlowFieldCFRSaleSg            = "c_fr_sale_sg"           // 销售商品、提供劳务收到的现金
	CashFlowFieldCInfFrOperateA       = "c_inf_fr_operate_a"     // 经营活动现金流入小计
	CashFlowFieldNCashflowAct         = "n_cashflow_act"         // 经营活动产生的现金流量净额
	CashFlowFieldStotInflowsInvAct    = "stot_inflows_inv_act"   // 投资活动现金流入小计
	CashFlowFieldNCashflowInvAct      = "n_cashflow_inv_act"     // 投资活动产生的现金流量净额
	CashFlowFieldStotCashInFncAct     = "stot_cash_in_fnc_act"   // 筹资活动现金流入小计
	CashFlowFieldNCashFlowsFncAct     = "n_cash_flows_fnc_act"   // 筹资活动产生的现金流量净额
	CashFlowFieldFreeCashflow         = "free_cashflow"          // 企业自由现金流量
	CashFlowFieldNIncrCashCashEqu     = "n_incr_cash_cash_equ"   // 现金及现金等价物净增加额
)

// CashFlowParams 现金流量表参数
// 接口: cashflow
// 描述: 获取上市公司现金流量表
// 文档: https://tushare.pro/document/2?doc_id=44
type CashFlowParams struct {
	TSCode     string   // 股票代码（必填）
	AnnDate    string   // 公告日期（YYYYMMDD格式）
	FAnnDate   string   // 实际公告日期
	StartDate  string   // 公告日开始日期
	EndDate    string   // 公告日结束日期
	Period     string   // 报告期
	ReportType string   // 报告类型
	CompType   CompType // 公司类型
	IsCalc     int      // 是否计算报表
	Fields     []string // 返回字段列表
}

// CashFlowItem 现金流量表响应项
type CashFlowItem struct {
	TSCode               string  `json:"ts_code"`                // TS股票代码
	AnnDate              string  `json:"ann_date"`               // 公告日期
	FAnnDate             string  `json:"f_ann_date"`             // 实际公告日期
	EndDate              string  `json:"end_date"`               // 报告期
	CompType             string  `json:"comp_type"`              // 公司类型
	ReportType           string  `json:"report_type"`            // 报表类型
	NetProfit            float64 `json:"net_profit"`             // 净利润
	FinanExp             float64 `json:"finan_exp"`              // 财务费用
	CFRSaleSg            float64 `json:"c_fr_sale_sg"`           // 销售商品、提供劳务收到的现金
	CInfFrOperateA       float64 `json:"c_inf_fr_operate_a"`     // 经营活动现金流入小计
	NCashflowAct         float64 `json:"n_cashflow_act"`         // 经营活动产生的现金流量净额
	StotInflowsInvAct    float64 `json:"stot_inflows_inv_act"`   // 投资活动现金流入小计
	NCashflowInvAct      float64 `json:"n_cashflow_inv_act"`     // 投资活动产生的现金流量净额
	StotCashInFncAct     float64 `json:"stot_cash_in_fnc_act"`   // 筹资活动现金流入小计
	NCashFlowsFncAct     float64 `json:"n_cash_flows_fnc_act"`   // 筹资活动产生的现金流量净额
	FreeCashflow         float64 `json:"free_cashflow"`          // 企业自由现金流量
	NIncrCashCashEqu     float64 `json:"n_incr_cash_cash_equ"`   // 现金及现金等价物净增加额
}

// CashFlow 获取现金流量表数据（自动处理分页）
func CashFlow(c *tushare.Client, params *CashFlowParams, opts ...tushare.QueryOption) ([]*CashFlowItem, error) {
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
	if params.IsCalc != 0 {
		reqParams["is_calc"] = params.IsCalc
	}

	fields := ""
	if len(params.Fields) > 0 {
		fields = strings.Join(params.Fields, ",")
	}

	resp, err := c.Query("cashflow", reqParams, fields, opts...)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, &tushare.APIError{
			Code: resp.Code,
			Msg:  resp.Msg,
		}
	}

	var items []*CashFlowItem
	if err := resp.ToStruct(&items); err != nil {
		return nil, err
	}

	return items, nil
}
