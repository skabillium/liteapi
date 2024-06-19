package main

import (
	"errors"
	"os"
	"skabillium/liteapi/api/health"
	"skabillium/liteapi/api/hotels"
	"skabillium/liteapi/clients"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
		return nil, errors.New("not HotelBeds url provided")
	}
	if !strings.HasPrefix(env.HotelBedsUrl, "http") {
		return nil, errors.New("invalid format for HotelBeds url")
	}

	if env.HotelBedsApiKey == "" {
		return nil, errors.New("not HotelBeds api key provided")
	}
	if env.HotelBedsSecret == "" {
		return nil, errors.New("not HotelBeds secret provided")
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

	r.GET("/health/status", health.GetStatusHandler)
	r.GET("/hotels", hotelsHanlders.GetHotelsHandler)

	r.Run()
}
