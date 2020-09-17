package analyzer

import (
	"Fachoi_fund_analysis/db_model"
	"Fachoi_fund_analysis/metrics"
	"Fachoi_fund_analysis/resource_manager"
	"Fachoi_fund_analysis/saver"
	"Fachoi_fund_analysis/util"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
)

type AnnualAnalyzer struct {
	codes                     []string
	num                       int
	db                        *sqlx.DB
	analyzeCount              chan int
	annualReturnsSaver        *saver.AnnualReturnsSaver
	annualAnalysisResultSaver *saver.AnnualAnalysisResultSaver
	metrics                   *metrics.AnnualMetrics
}

func NewAnnualAnalyzer(codes []string, num int, db *sqlx.DB, filename string, metrics *metrics.AnnualMetrics) *AnnualAnalyzer {
	return &AnnualAnalyzer{
		codes,
		num,
		db,
		make(chan int, len(codes)),
		saver.NewAnnualReturnsSaver(db),
		saver.NewAnnualAnalysisResultSaver(filename),
		metrics,
	}
}

// 计算并存储所有基金的历史年化回报率，为下一步分析做预处理
func (aa *AnnualAnalyzer) Preprocess() {
	arm := resource_manager.NewResourceManager(aa.num)
	for _, code := range aa.codes {
		arm.GetOne()
		go func(code string) {
			defer arm.FreeOne()
			aa.process(code)
		}(code)
	}
}

func (aa *AnnualAnalyzer) process(code string) {
	util.TruncateTable("annual_returns_table", aa.db)

	startDay, historyData, has := util.GetHistoryDataByCode(code, aa.db)
	if has == false || len(historyData) == 1 {
		return
	}

	lastTradeDays := util.GetLastTradeDays(startDay, historyData) // 每年最后一个交易日
	annualReturnRatios := make([]float32, 0)                      // 年化回报
	years := make([]int, 0)                                       // 年份
	for i, lastTradeDay := range lastTradeDays {
		if i == 0 {
			annualReturnRatios = append(annualReturnRatios, util.CalculateReturnRatio(historyData[lastTradeDay], historyData[startDay]))
		} else {
			annualReturnRatios = append(annualReturnRatios, util.CalculateReturnRatio(historyData[lastTradeDay], historyData[lastTradeDays[i-1]]))
		}
		years = append(years, time.Unix(lastTradeDay, 0).Year())
	}
	arm := new(db_model.AnnualReturnsModel)
	arm.Code = code
	arm.Years = util.SliceToString(years)
	arm.AnnualReturns = util.SliceToString(annualReturnRatios)
	aa.annualReturnsSaver.Save(arm)

	aa.analyzeCount <- 1
	fmt.Printf("年化回报率预处理进度：%d / %d\n", len(aa.analyzeCount), cap(aa.analyzeCount))
}

func (aa *AnnualAnalyzer) Analyze() {
	LaunchYears := aa.metrics.LaunchYears
	NegativeReturnsYears := aa.metrics.NegativeReturnsYears
	selectedFundTypes := aa.metrics.SelectedFundTypes
	candidateCodes := make([]string, 0)
	arm := new(db_model.AnnualReturnsModel)
	fim := new(db_model.FundInfoModel)

	sqlStr := "SELECT fund_code,years,annual_returns FROM annual_returns_table"
	rows, _ := aa.db.Query(sqlStr)
	for rows.Next() {
		rows.Scan(&arm.Code, &arm.Years, &arm.AnnualReturns)
		annualReturnsStr := strings.ReplaceAll(arm.AnnualReturns, "[", "")
		annualReturnsStr = strings.ReplaceAll(annualReturnsStr, "]", "")
		annualReturns := strings.Split(annualReturnsStr, ",")
		if len(annualReturns) <= LaunchYears {
			// 跳过成立时间小于LaunchYears的基金
			continue
		}
		negativeReturnsCount := 0
		for _, annualReturnsStr := range annualReturns {
			annualReturnFloat, _ := strconv.ParseFloat(annualReturnsStr, 32)
			if annualReturnFloat < 0 {
				negativeReturnsCount++
			}
		}
		if negativeReturnsCount <= NegativeReturnsYears {
			// 负收益年数小于NegativeReturnsYears的基金加入候选数组
			candidateCodes = append(candidateCodes, arm.Code)
		}
	}
	rows.Close()

	result := ""
	sqlStr = "SELECT fund_code_front_end,fund_short_name,fund_type,fund_asset_size FROM fund_info_table WHERE ("
	for i, code := range candidateCodes {
		if i == len(candidateCodes)-1 {
			sqlStr += "fund_code_front_end=" + code + ") AND ("
		} else {
			sqlStr += "fund_code_front_end=" + code + " OR "
		}
	}

	for i, types := range selectedFundTypes {
		if i == len(selectedFundTypes)-1 {
			sqlStr += "fund_type='" + types + "')"
		} else {
			sqlStr += "fund_type='" + types + "' OR "
		}
	}
	rows, _ = aa.db.Query(sqlStr)
	for rows.Next() {
		rows.Scan(&fim.CodeFront, &fim.ShortName, &fim.FundType, &fim.AssetSize)
		result += fim.CodeFront + " " + fim.ShortName + " " + fim.FundType + " " + fmt.Sprintf("%f", fim.AssetSize) + "\n"

	}
	rows.Close()
	aa.annualAnalysisResultSaver.Save(result, aa.metrics)
}
