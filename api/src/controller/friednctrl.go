package controller

import (
	"bean"
	"dao"
	"fmt"
)

func GetFriends(userId int64) *bean.ResInfo {

	friends := dao.GetFriends(userId)
	res := &bean.ResInfo{Code: "000001"}
	res.Desc = "success"
	res.Data = friends
	fmt.Println(friends)
	return res
}
