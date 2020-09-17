package db_model

type AnnualReturnsModel struct {
	Id            int    `db:"id"`
	Code          string `db:"fund_code"`
	Years         string `db:"years"`
	AnnualReturns string `db:"annual_returns"`
}
