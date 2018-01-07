package bean

import (
	"time"
)

/**
 des  用户信息结构体
 author liupengfei

**/
type User struct {
	UserId      int64     `json:"userId"`      //'业务主键',
	NickName    string    `json:"nickName"`    // '昵称',
	MobilePhone string    `json:"mobilePhone"` // '手机号',
	Password    string    `json:"password"`    // '密码',
	Salt        string    `json:"salt"`        // '密码盐值',
	Email       string    `json:"email"`       // '用户邮箱',
	Username    string    `json:"username"`    // '用户名',
	Sex         int8      `json:"sex"`         // '性别',
	Birthday    time.Time `json:"birthday"`    // '出生日期',
	SyncKey     int64     `json:"syncKey"`     // '消息同步序号',
	CreateTime  time.Time `json:"createTime"`  //'创建时间',
	ModifyTime  time.Time `json:"modifyTime"`  // '更新时间',
	Token       string    `json:"token"`       // '用户登录凭证',
}
