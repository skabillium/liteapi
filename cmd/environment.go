package main

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const DefaultPort = "8080"

type Environment struct {
	GoEnv           string
	Port            int
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
	portStr := cmp.Or(os.Getenv("PORT"), DefaultPort)
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		return nil, fmt.Errorf("invalid port '%s', port needs to be a positive integer", portStr)
	}

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
		return nil, fmt.Errorf("invalid format for HotelBeds url: %s", env.HotelBedsUrl)
	}

	if env.HotelBedsApiKey == "" {
		return nil, errors.New("no HotelBeds api key provided")
	}
	if env.HotelBedsSecret == "" {
		return nil, errors.New("no HotelBeds secret provided")
	}

	return env, nil
}
