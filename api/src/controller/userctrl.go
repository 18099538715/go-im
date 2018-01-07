package controller

import (
	"bean"
	"crypto/md5"
	"dao"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/valyala/fasthttp"
)

func UserRegister(ctx *fasthttp.RequestCtx) {
	body := ctx.Request.Body() //获取post的数据
	res := &bean.ResInfo{}
	user := &bean.User{}
	err := json.Unmarshal(body, user)
	if err != nil {
		fmt.Println(err)
	}
	user.Salt = GetRandomString()
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(user.Password + user.Salt))
	cipherStr := md5Ctx.Sum(nil)
	user.Password = hex.EncodeToString(cipherStr)
	fmt.Print(user.Password)
	dao.RegisterUser(user)
	res.Code = "000001"
	res.Desc = "success"
	b, err := json.Marshal(res)
	ctx.Response.SetBody(b)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func UserLogin(ctx *fasthttp.RequestCtx) {
	user := &bean.User{UserId: 1, Token: "aaaa", MobilePhone: "20214305"}
	res := &bean.ResInfo{Code: "000001"}
	res.Data = user
	b, _ := json.Marshal(res)
	ctx.Response.SetBody(b)
	ctx.SetStatusCode(fasthttp.StatusOK)
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
