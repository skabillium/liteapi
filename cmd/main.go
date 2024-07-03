package main

import (
	"fmt"
	"skabillium/liteapi/cmd/api/health"
	"skabillium/liteapi/cmd/api/hotels"
	"skabillium/liteapi/cmd/clients"
	"time"

	"github.com/gin-gonic/gin"

	_ "skabillium/liteapi/cmd/api/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	env, err := loadEnv()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	hbc := clients.NewHotelBedsClient(env.HotelBedsApiKey, env.HotelBedsSecret, env.HotelBedsUrl, 10*time.Second)
	hotels.RegisterHotelHandlers(r, hbc)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/health/status", health.GetStatusHandler)

	r.Run(fmt.Sprintf(":%d", env.Port))
}
