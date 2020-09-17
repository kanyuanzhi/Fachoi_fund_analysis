package main

import (
	"Fachoi_fund_analysis/analyzer"
	"Fachoi_fund_analysis/db_mysql"
	"Fachoi_fund_analysis/metrics"
	"Fachoi_fund_analysis/util"
)

func main() {
	mysqlDB := db_mysql.NewMysql()
	mysqlDB.InitDatabase()
	db := mysqlDB.GetDB()
	fundCodes := util.GetFundCodes(db)

	filename := "annual_result.txt"        // 分析结果存储文件名，将保存到analysis_results目录中
	launchYears := 5                       // 基金成立时长（年）大于>=launchYears
	negativeReturnsYears := 2              // 基金在在选定时长内年化负收益的年份小于等于negativeReturnsYears的次数
	selectedFundTypesInInt := []int{9, 18} // 所要分析的基金类型，数字对应关系详见metrics/annual_metrics
	annualMetrics := metrics.NewAnnualMetrics(launchYears, negativeReturnsYears, selectedFundTypesInInt)
	aa := analyzer.NewAnnualAnalyzer(fundCodes, 20, db, filename, annualMetrics)
	//aa.Preprocess()
	aa.Analyze()
}
