package util

import (
	"Fachoi_fund_analysis/db_model"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func GetFundCodes(db *sqlx.DB) []string {
	flm := new(db_model.FundListModel)
	fundCodes := make([]string, 0)
	sqlStr := "SELECT fund_code FROM fund_list_table"
	rows, _ := db.Query(sqlStr)
	for rows.Next() {
		rows.Scan(&flm.Code)
		fundCodes = append(fundCodes, flm.Code)
	}
	rows.Close()
	return fundCodes
}

func GetHistoryDataByCode(code string, db *sqlx.DB) (int64, map[int64]float32, bool) {
	fhm := new(db_model.FundHistoryModel)
	historyData := make(map[int64]float32)
	var startDay int64

	tableName := "history_" + code + "_table"
	sqlStr := "SELECT date, accumulated_net_asset_value FROM " + tableName + " LIMIT 1"

	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Println(err)
		return 0, nil, false
	}
	for rows.Next() {
		rows.Scan(&fhm.Date, &fhm.AccumulatedValue)
	}
	rows.Close()

	if fhm.AccumulatedValue == 0 {
		return 0, nil, false
	}

	startDay = fhm.Date

	sqlStr = "SELECT date, accumulated_net_asset_value FROM " + tableName
	rows, _ = db.Query(sqlStr)
	for rows.Next() {
		rows.Scan(&fhm.Date, &fhm.AccumulatedValue)
		historyData[fhm.Date] = fhm.AccumulatedValue
	}
	rows.Close()
	return startDay, historyData, true
}

func TruncateTable(tableName string, db *sqlx.DB) {
	sqlStr := "TRUNCATE TABLE " + tableName
	_, err := db.Exec(sqlStr)
	CheckError(err, "TruncateTable")
}

func CreateFundHistoryTable(code string, db *sqlx.DB) int64 {
	tableName := "history_" + code + "_table"
	sqlStr := "CREATE TABLE IF NOT EXISTS " + tableName + " (" +
		"id INT AUTO_INCREMENT," +
		"date BIGINT," +
		"date_string VARCHAR(50)," +
		"net_asset_value FLOAT, " +
		"accumulated_net_asset_value FLOAT, " +
		"earnings_per_10000 FLOAT, " +
		"7_day_annual_return FLOAT, " +
		"PRIMARY KEY (id)" +
		")"
	result, err := db.Exec(sqlStr)

	CheckError(err, "createFundInfoTable")
	rowNum, _ := result.RowsAffected()
	return rowNum
}
