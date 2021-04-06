package router

import (
	"github.com/whuangz/go-example/go-api/handler"
	"github.com/whuangz/go-example/go-api/repository"
	"github.com/whuangz/go-example/go-api/service"
)

func blogRoutes() {
	repo := repository.NewPostRepo(database)
	service := service.NewPostService(repo)
	handler.NewPostHandler(router, service)
}
