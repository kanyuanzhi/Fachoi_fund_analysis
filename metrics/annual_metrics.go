package metrics

type AnnualMetrics struct {
	LaunchYears          int      // 成立时间大于等于years
	NegativeReturnsYears int      // 负收益年数
	SelectedFundTypes    []string // 所选择的基金类型
}

func NewAnnualMetrics(launchYears int, negativeReturnsYears int, selectedFundTypesInInt []int) *AnnualMetrics {
	allFundTypes := AllFundTypes()
	selectedFundTypesInStr := make([]string, 0)
	for _, t := range selectedFundTypesInInt {
		selectedFundTypesInStr = append(selectedFundTypesInStr, allFundTypes[t])
	}
	return &AnnualMetrics{
		launchYears,
		negativeReturnsYears,
		selectedFundTypesInStr,
	}
}

func AllFundTypes() map[int]string {
	allFundTypes := make(map[int]string, 19)
	allFundTypes[1] = "混合型"
	allFundTypes[2] = "债券型"
	allFundTypes[3] = "定开债券"
	allFundTypes[4] = "联接基金"
	allFundTypes[5] = "货币型"
	allFundTypes[6] = "QDII"
	allFundTypes[7] = "股票指数"
	allFundTypes[8] = "QDII-指数"
	allFundTypes[9] = "股票型"
	allFundTypes[10] = "理财型"
	allFundTypes[11] = "债券指数"
	allFundTypes[12] = "保本型"
	allFundTypes[13] = "其他创新"
	allFundTypes[14] = "混合-FOF"
	allFundTypes[15] = "股票-FOF"
	allFundTypes[16] = "固定收益"
	allFundTypes[17] = "分级杠杆"
	allFundTypes[18] = "ETF-场内"
	allFundTypes[19] = "QDII-ETF"
	return allFundTypes
}
