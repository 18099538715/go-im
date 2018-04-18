package dao

import (
	"bean"
	"dbcon"
	"fmt"
)

func GetFriends(userId int64) []*bean.Friend {
	friends := []*bean.Friend{}
	stmt, err := dbcon.MySqlCon.Prepare("select friendId,userId,friendName from t_friend where userId=?")
	rows, err := stmt.Query(userId)
	if err != nil {
		fmt.Printf("获取好友信息出错: %v\n", err)
		return friends
	}
	for rows.Next() {
		friend := &bean.Friend{}
		rows.Scan(&friend.FriendId, &friend.UserId, &friend.FriendName)
		friends = append(friends, friend)
	}
	return friends
}
