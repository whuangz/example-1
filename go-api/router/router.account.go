package router

import (
	"github.com/whuangz/go-example/go-api/handler"
	"github.com/whuangz/go-example/go-api/service"
)

func accountRoutes() {
	service := service.NewAccountService(nil)
	handler.NewAccountHandler(router, service)
}
