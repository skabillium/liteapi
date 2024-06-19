package main

import (
	"cmp"
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const DefaultPort = "8080"

type Environment struct {
	GoEnv           string
	Port            string
	HotelBedsUrl    string
	HotelBedsApiKey string
	HotelBedsSecret string
}

// Reads enviroment variables, validates them and returns an Envrironment struct
func loadEnv() (*Environment, error) {
	// Default environment is "development"
	goEnv := cmp.Or(os.Getenv("GO_ENV"), "development")
	if goEnv != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}
	port := cmp.Or(os.Getenv("PORT"), DefaultPort)

	// Validate environment
	env := &Environment{
		GoEnv:           goEnv,
		Port:            port,
		HotelBedsUrl:    os.Getenv("HB_URL"),
		HotelBedsApiKey: os.Getenv("HB_API_KEY"),
		HotelBedsSecret: os.Getenv("HB_SECRET"),
	}

	// Validate environment variables
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
