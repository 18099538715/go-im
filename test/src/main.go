package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/18099538715/go-rest/webrestful"
)

type User struct {
	Username string `json:"userName"`
	Password string `json:"password"`
}

type UserController struct {
}

func add(user User) User {
	return user
}

func main() {
	handler := webrestful.Handler{}
	webrestful.Route("/aaa/{userId}", http.MethodPost, "application/json", func(userId *string, user User) User {
		b, _ := json.Marshal(user)
		fmt.Println(string(b), *userId)
		return user
	})
	webrestful.Route("/aaa/{userId}/{uesrname}", http.MethodGet, "application/json", func(userId string, userName string) string {
		return userId
	})
	http.ListenAndServe(":8000", handler)
}
