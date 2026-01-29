package financial_test

import (
	"fmt"
	"log"

	tushare "github.com/fletcherlau/go-tushare"
	"github.com/fletcherlau/go-tushare/stock/financial"
)

func ExampleIncome() {
	// 创建客户端
	client := tushare.NewClient("your_token")

	// 获取平安银行的利润表数据
	items, err := financial.Income(client, &financial.IncomeParams{
		TSCode: "000001.SZ",
		Period: "20241231",
		Fields: []string{
			financial.IncomeFieldAnnDate,
			financial.IncomeFieldEndDate,
			financial.IncomeFieldBasicEps,
			financial.IncomeFieldTotalRevenue,
			financial.IncomeFieldNIncome,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("获取 %d 条记录\n", len(items))
	if len(items) > 0 {
		fmt.Printf("第一条: 报告期=%s 基本每股收益=%.4f 净利润=%.2f\n",
			items[0].EndDate, items[0].BasicEps, items[0].NIncome)
	}
}

func ExampleBalanceSheet() {
	// 创建客户端
	client := tushare.NewClient("your_token")

	// 获取平安银行的资产负债表数据
	items, err := financial.BalanceSheet(client, &financial.BalanceSheetParams{
		TSCode: "000001.SZ",
		Period: "20241231",
		Fields: []string{
			financial.BalanceSheetFieldAnnDate,
			financial.BalanceSheetFieldEndDate,
			financial.BalanceSheetFieldTotalAssets,
			financial.BalanceSheetFieldTotalLiab,
			financial.BalanceSheetFieldTotalHldrEqyIncMinInt,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("获取 %d 条记录\n", len(items))
	if len(items) > 0 {
		fmt.Printf("第一条: 报告期=%s 总资产=%.2f 总负债=%.2f 股东权益=%.2f\n",
			items[0].EndDate, items[0].TotalAssets, items[0].TotalLiab, items[0].TotalHldrEqyIncMinInt)
	}
}

func ExampleCashFlow() {
	// 创建客户端
	client := tushare.NewClient("your_token")

	// 获取平安银行的现金流量表数据
	items, err := financial.CashFlow(client, &financial.CashFlowParams{
		TSCode: "000001.SZ",
		Period: "20241231",
		Fields: []string{
			financial.CashFlowFieldAnnDate,
			financial.CashFlowFieldEndDate,
			financial.CashFlowFieldNCashflowAct,
			financial.CashFlowFieldNCashflowInvAct,
			financial.CashFlowFieldNCashFlowsFncAct,
			financial.CashFlowFieldFreeCashflow,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("获取 %d 条记录\n", len(items))
	if len(items) > 0 {
		fmt.Printf("第一条: 报告期=%s 经营现金流=%.2f 投资现金流=%.2f 筹资现金流=%.2f\n",
			items[0].EndDate, items[0].NCashflowAct, items[0].NCashflowInvAct, items[0].NCashFlowsFncAct)
	}
}

func ExampleFinaIndicator() {
	// 创建客户端
	client := tushare.NewClient("your_token")

	// 获取平安银行的财务指标数据
	items, err := financial.FinaIndicator(client, &financial.FinaIndicatorParams{
		TSCode: "000001.SZ",
		Period: "20241231",
		Fields: []string{
			financial.FinaIndicatorFieldAnnDate,
			financial.FinaIndicatorFieldEndDate,
			financial.FinaIndicatorFieldRoe,
			financial.FinaIndicatorFieldRoa,
			financial.FinaIndicatorFieldDebtToAssets,
			financial.FinaIndicatorFieldProfitDedt,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("获取 %d 条记录\n", len(items))
	if len(items) > 0 {
		fmt.Printf("第一条: 报告期=%s ROE=%.4f ROA=%.4f 资产负债率=%.4f\n",
			items[0].EndDate, items[0].Roe, items[0].Roa, items[0].DebtToAssets)
	}
}
