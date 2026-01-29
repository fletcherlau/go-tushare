package financial

import (
	"strings"

	tushare "github.com/fletcherlau/go-tushare"
)

// BalanceSheetField 返回字段常量
const (
	BalanceSheetFieldTSCode                = "ts_code"                   // TS股票代码
	BalanceSheetFieldAnnDate               = "ann_date"                  // 公告日期
	BalanceSheetFieldFAnnDate              = "f_ann_date"                // 实际公告日期
	BalanceSheetFieldEndDate               = "end_date"                  // 报告期
	BalanceSheetFieldReportType            = "report_type"               // 报表类型
	BalanceSheetFieldCompType              = "comp_type"                 // 公司类型
	BalanceSheetFieldTotalShare            = "total_share"               // 期末总股本
	BalanceSheetFieldCapRese               = "cap_rese"                  // 资本公积金
	BalanceSheetFieldUndistrPorfit         = "undistr_porfit"            // 未分配利润
	BalanceSheetFieldMoneyCap              = "money_cap"                 // 货币资金
	BalanceSheetFieldTotalAssets           = "total_assets"              // 资产总计
	BalanceSheetFieldTotalLiab             = "total_liab"                // 负债合计
	BalanceSheetFieldTotalHldrEqyIncMinInt = "total_hldr_eqy_inc_min_int" // 股东权益合计（含少数股东权益）
	BalanceSheetFieldTotalLiabHldrEqy      = "total_liab_hldr_eqy"       // 负债及股东权益总计
)

// BalanceSheetParams 资产负债表参数
// 接口: balancesheet
// 描述: 获取上市公司资产负债表
// 文档: https://tushare.pro/document/2?doc_id=36
type BalanceSheetParams struct {
	TSCode     string   // 股票代码（必填）
	AnnDate    string   // 公告日期（YYYYMMDD格式）
	StartDate  string   // 公告日开始日期
	EndDate    string   // 公告日结束日期
	Period     string   // 报告期（每个季度最后一天的日期）
	ReportType string   // 报告类型
	CompType   CompType // 公司类型：1一般工商业 2银行 3保险 4证券
	Fields     []string // 返回字段列表
}

// BalanceSheetItem 资产负债表响应项
type BalanceSheetItem struct {
	TSCode                string  `json:"ts_code"`                  // TS股票代码
	AnnDate               string  `json:"ann_date"`                 // 公告日期
	FAnnDate              string  `json:"f_ann_date"`               // 实际公告日期
	EndDate               string  `json:"end_date"`                 // 报告期
	ReportType            string  `json:"report_type"`              // 报表类型
	CompType              string  `json:"comp_type"`                // 公司类型
	TotalShare            float64 `json:"total_share"`              // 期末总股本
	CapRese               float64 `json:"cap_rese"`                 // 资本公积金
	UndistrPorfit         float64 `json:"undistr_porfit"`           // 未分配利润
	MoneyCap              float64 `json:"money_cap"`                // 货币资金
	TotalAssets           float64 `json:"total_assets"`             // 资产总计
	TotalLiab             float64 `json:"total_liab"`               // 负债合计
	TotalHldrEqyIncMinInt float64 `json:"total_hldr_eqy_inc_min_int"` // 股东权益合计（含少数股东权益）
	TotalLiabHldrEqy      float64 `json:"total_liab_hldr_eqy"`      // 负债及股东权益总计
}

// BalanceSheet 获取资产负债表数据（自动处理分页）
func BalanceSheet(c *tushare.Client, params *BalanceSheetParams, opts ...tushare.QueryOption) ([]*BalanceSheetItem, error) {
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

	resp, err := c.Query("balancesheet", reqParams, fields, opts...)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, &tushare.APIError{
			Code: resp.Code,
			Msg:  resp.Msg,
		}
	}

	var items []*BalanceSheetItem
	if err := resp.ToStruct(&items); err != nil {
		return nil, err
	}

	return items, nil
}
