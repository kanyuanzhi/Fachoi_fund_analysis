package db_mysql

import (
	"Fachoi_fund_analysis/util"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type MysqlDB struct {
	db *sqlx.DB // 数据库连接池
}

func NewMysql() *MysqlDB {
	return &MysqlDB{}
}

// 初始化数据库，包括建立数据库，建立数据表
func (m *MysqlDB) InitDatabase() {
	m.db = createConnection()
	createAnnualizedReturnTable(m.db)
}

// 获取连接
func (m *MysqlDB) GetDB() *sqlx.DB {
	if m.db == nil {
		fmt.Println("createConnection")
		m.db = createConnection()
	}
	return m.db
}

func createConnection() *sqlx.DB {
	user, pass, host, port, dbname, charset := util.GetDBConfig()
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", user, pass, host, port, dbname, charset)
	db, err := sqlx.Connect("mysql", dbDSN)
	if err != nil {
		log.Println("dbDSN: " + dbDSN)
		panic("数据源配置不正确: " + err.Error())
	}
	// 最大连接数
	db.SetMaxOpenConns(100)
	// 闲置连接数
	db.SetMaxIdleConns(20)
	// 最大连接周期
	db.SetConnMaxLifetime(120 * time.Second)

	if err = db.Ping(); nil != err {
		panic("数据库链接失败: " + err.Error())
	}
	return db
}

func createAnnualizedReturnTable(db *sqlx.DB) {
	sqlStr := "CREATE TABLE IF NOT EXISTS annual_returns_table (" +
		"id INT AUTO_INCREMENT, " +
		"fund_code VARCHAR(10), " +
		"years VARCHAR(255), " +
		"annual_returns VARCHAR(255), " +
		"PRIMARY KEY (id))"
	_, err := db.Exec(sqlStr)
	util.CheckError(err, "createFundListTable")
}
