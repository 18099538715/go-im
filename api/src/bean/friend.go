package bean

/**
 des  好友信息
 author liupengfei

**/
type Friend struct {
	UserId     int64  `json:"userId"`     //'业务主键',
	FriendName string `json:"friendName"` // '好友姓名',
	FriendId   int64  `json:"friendId"`   //'业务主键',
}
