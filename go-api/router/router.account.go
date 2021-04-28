package router

import (
	"github.com/whuangz/go-example/go-api/config"
	"github.com/whuangz/go-example/go-api/handler"
	"github.com/whuangz/go-example/go-api/repository"
	"github.com/whuangz/go-example/go-api/service"
)

func accountRoutes() {

	accRepo := repository.NewAccountRepo(database)
	accService := service.NewAccountService(accRepo)
	tokenRepo := repository.NewTokenRepo(redisClient)
	tokenService := service.NewTokenService(
		tokenRepo,
		config.JWT_PRIVATE_KEY, config.JWT_PUBLIC_KEY, config.JWT_REFRESH_SECRET,
		config.ACCESS_TOKEN_EXP, config.REFRESH_TOKEN_EXP)

	handler.NewAccountHandler(router, accService, tokenService)
}
