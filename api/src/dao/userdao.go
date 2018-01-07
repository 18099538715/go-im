package dao

import (
	"bean"
	"database/sql"
	"fmt"
)

func RegisterUser(user *bean.User) {
	stmt, err := db.Prepare("INSERT INTO user(mobilePhone, password,salt,syncKey) VALUES(?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stmt.Exec(user.MobilePhone, user.Password, user.Salt, 0)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

}
func GetUser() {

	rows, err := db.Query("select * from user")
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}
