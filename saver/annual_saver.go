package saver

import (
	"Fachoi_fund_analysis/db_model"
	"Fachoi_fund_analysis/metrics"
	"Fachoi_fund_analysis/util"
	"bufio"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"path"
	"strings"
)

// save preprocessed data to mysql
type AnnualReturnsSaver struct {
	db *sqlx.DB
}

func NewAnnualReturnsSaver(db *sqlx.DB) *AnnualReturnsSaver {
	return &AnnualReturnsSaver{
		db,
	}
}

func (ars *AnnualReturnsSaver) Save(arm *db_model.AnnualReturnsModel) {
	sqlStr := "INSERT INTO annual_returns_table(fund_code,years,annual_returns) VALUES(?,?,?)"
	smtp, _ := ars.db.Prepare(sqlStr)
	_, err := smtp.Exec(arm.Code, arm.Years, arm.AnnualReturns)
	util.CheckError(err, "save AnnualReturnModel")
}

// save final analysis result to txt
type AnnualAnalysisResultSaver struct {
	filename string
}

func NewAnnualAnalysisResultSaver(filename string) *AnnualAnalysisResultSaver {
	return &AnnualAnalysisResultSaver{
		filename,
	}
}

func (aars *AnnualAnalysisResultSaver) Save(result string, metrics *metrics.AnnualMetrics) {
	filePath := path.Join("analysis_results", aars.filename)
	pathExists, err := util.PathExists(filePath)
	util.CheckError(err, "PathExists")
	if pathExists == false {
		_, err = os.Create(filePath)
		util.CheckError(err, "AnnualAnalysisResultSaver Save Create")
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}
	defer file.Close()
	dstStr := "METRICS: %d年内年化收益为负的次数少于等于%d次的基金（"
	dstTypesStr := strings.Join(metrics.SelectedFundTypes, ",")
	dstStr += dstTypesStr + "）\n" + result + "\n"
	writer := bufio.NewWriter(file)
	writer.WriteString(fmt.Sprintf(dstStr, metrics.LaunchYears, metrics.NegativeReturnsYears))
	writer.Flush()
}
