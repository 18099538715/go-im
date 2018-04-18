package dao

import (
	"bean"
	"dbcon"
	"fmt"
)

func RegisterUser(user *bean.User) {
	stmt, err := dbcon.MySqlCon.Prepare("INSERT INTO t_user(mobilePhone, password,salt,syncKey) VALUES(?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stmt.Exec(user.MobilePhone, user.Password, user.Salt, 0)
	if err != nil {
		fmt.Println(err)
	}

}
func GetUser(mobilePhone string) *bean.User {
	stmt, err := dbcon.MySqlCon.Prepare("select password,salt,userId from t_user where mobilePhone=?")
	user := &bean.User{}
	err = stmt.QueryRow(mobilePhone).Scan(&user.Password, &user.Salt, &user.UserId)
	if err != nil {
		fmt.Println(err)
	}
	return user
}
