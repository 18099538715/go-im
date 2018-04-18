package main

import (
	"controller"
	"net/http"

	"github.com/18099538715/go-rest/webrestful"
)

func main() {

	handler := webrestful.Handler{}
	webrestful.Route("/user/register", http.MethodPost, "application/json", controller.UserRegister)
	webrestful.Route("/user/login", http.MethodPost, "application/json", controller.Login)
	webrestful.Route("/friends/{userId}", http.MethodGet, "application/json", controller.GetFriends)
	http.ListenAndServe(":8099", handler)
}
