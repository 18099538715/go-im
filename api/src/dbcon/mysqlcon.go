package dbcon

import (
	"database/sql"
	"fmt"
)
import _ "github.com/go-sql-driver/mysql"

var MySqlCon *sql.DB

func init() {
	var err error
	MySqlCon, err = sql.Open("mysql", "root:123456@tcp(localhost:3306)/uchat?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}
	MySqlCon.SetMaxOpenConns(2000)
	MySqlCon.SetMaxIdleConns(1000)
	err = MySqlCon.Ping()
	if err != nil {
		fmt.Println(err)
	}
}
