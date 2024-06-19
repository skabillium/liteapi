package main

import (
	"errors"
	"os"
	"skabillium/liteapi/cmd/api/health"
	"skabillium/liteapi/cmd/api/hotels"
	"skabillium/liteapi/cmd/clients"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "skabillium/liteapi/cmd/api/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Environment struct {
	HotelBedsUrl    string
	HotelBedsApiKey string
	HotelBedsSecret string
}

func readEnv() (*Environment, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// Validate environment
	env := &Environment{
		HotelBedsUrl:    os.Getenv("HB_URL"),
		HotelBedsApiKey: os.Getenv("HB_API_KEY"),
		HotelBedsSecret: os.Getenv("HB_SECRET"),
	}

	if env.HotelBedsUrl == "" {
		return nil, errors.New("no HotelBeds url provided")
	}
	if !strings.HasPrefix(env.HotelBedsUrl, "http") {
		return nil, errors.New("invalid format for HotelBeds url")
	}

	if env.HotelBedsApiKey == "" {
		return nil, errors.New("no HotelBeds api key provided")
	}
	if env.HotelBedsSecret == "" {
		return nil, errors.New("no HotelBeds secret provided")
	}

	return env, nil
}

func main() {
	env, err := readEnv()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	hbc := clients.NewHotelBedsClient(env.HotelBedsApiKey, env.HotelBedsSecret, env.HotelBedsUrl, 10*time.Second)

	hotelsHanlders := hotels.NewHotelsHanlders(hbc)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/health/status", health.GetStatusHandler)
	r.GET("/hotels", hotelsHanlders.GetHotelsHandler)

	r.Run()
}
