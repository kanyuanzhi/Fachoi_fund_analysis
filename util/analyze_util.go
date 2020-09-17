package util

import (
	"fmt"
	"strconv"
	"time"
)

// 得到每年最后一个交易日期，如果是今年则返回离今天最近的一个交易日
func GetLastTradeDays(start int64, historyData map[int64]float32) []int64 {
	startDay := time.Unix(start, 0)
	lastDayInThisYear := time.Date(startDay.Year(), 12, 31, 8, 0, 0, 0, startDay.Location())
	lastDays := make([]time.Time, 0)
	now := time.Now()
	for lastDayInThisYear.Before(now) {
		lastDays = append(lastDays, lastDayInThisYear)
		lastDayInThisYear = time.Date(lastDayInThisYear.Year()+1, 12, 31, 8, 0, 0, 0, lastDayInThisYear.Location())
	}

	// 确定历年最后一个交易日（即historyData中有当天的数据）
	lastTradeDays := make([]int64, 0)
	for _, lastDay := range lastDays {
		tempDay := lastDay.Unix()
		for {
			if _, has := historyData[tempDay]; has == true {
				lastTradeDays = append(lastTradeDays, tempDay)
				break
			} else {
				tempDay = tempDay - 3600*24 //减一天
			}
		}
	}

	//确定今年最后一个交易日
	latestLastTradeDay := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
	tempDay := latestLastTradeDay.Unix()
	for {
		if _, has := historyData[tempDay]; has == true {
			lastTradeDays = append(lastTradeDays, tempDay)
			break
		} else {
			tempDay = tempDay - 3600*24 //减一天
		}
	}
	return lastTradeDays
}

// 计算收益率，v1:较新的日期，v2:较旧的日期
func CalculateReturnRatio(v1 float32, v2 float32) float32 {
	returnRatio, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", (v1-v2)/v2), 32)
	return float32(returnRatio)
}
