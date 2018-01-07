package dao

import (
	"database/sql"
	"fmt"
)
import _ "github.com/go-sql-driver/mysql"

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(localhost:3306)/uchat?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
}
