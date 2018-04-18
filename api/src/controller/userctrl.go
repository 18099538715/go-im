package controller

import (
	"bean"
	"crypto/md5"
	"dao"
	"encoding/hex"
	"fmt"
	"math/rand"
	"rediscache"
	"strconv"
	"time"
)

func UserRegister(user *bean.User) *bean.ResInfo {
	user.Salt = GetRandomString()
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(user.Password + user.Salt))
	cipherStr := md5Ctx.Sum(nil)
	user.Password = hex.EncodeToString(cipherStr)
	dao.RegisterUser(user)
	res := &bean.ResInfo{Code: "000001"}
	res.Desc = "success"
	return res
}

func Login(loginUser *bean.User) *bean.ResInfo {
	fmt.Println(loginUser)
	res := &bean.ResInfo{Code: "000001"}
	user := dao.GetUser(loginUser.MobilePhone)
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(loginUser.Password + user.Salt))
	cipherStr := md5Ctx.Sum(nil)
	pwd := hex.EncodeToString(cipherStr)
	if user.Password == pwd {
		userReturn := &bean.User{}
		res.Desc = "success"
		userReturn.UserId = user.UserId
		userReturn.Token = getToken(user.UserId)
		res.Data = userReturn
		_, err := rediscache.SetUserToken(user.UserId, userReturn.Token)
		if err != nil {
			fmt.Println("登录保存到redis出错", err)
		}
	} else {
		res.Code = "000002"
		res.Desc = "error"
	}

	return res
}
func GetRandomString() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
func getToken(userId int64) string {
	crutime := time.Now().Unix()
	s := strconv.FormatInt(crutime, 10) + strconv.FormatInt(userId, 10)
	h := md5.New()
	h.Write([]byte(s))
	cipherStr := h.Sum(nil)
	token := hex.EncodeToString(cipherStr)
	fmt.Println("token--->", token)
	return token + strconv.FormatInt(crutime, 10)
}
