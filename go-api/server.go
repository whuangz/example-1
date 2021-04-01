package main

import (
	_ "github.com/whuangz/go-example/go-api/routers"
	"github.com/whuangz/go-example/go-api/engine"

)

func main() {
	engine.Connect()
}