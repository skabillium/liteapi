package clients

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Struct that holds all functionality specific to the HotelBeds API like:
// - Endpoints
// - Auth
// - Response parsing and validation
type HotelBedsClient struct {
	apiKey     string
	secret     string
	apiAddress string
	httpClient http.Client
}

// Instantiates a new HotelBedsClient, this is supposed to be called only once (on application startup).
// The client instance is meant to be shared across the server.
func NewHotelBedsClient(apikey, secret, apiAddress string, timeout time.Duration) *HotelBedsClient {
	return &HotelBedsClient{
		apiKey:     apikey,
		secret:     secret,
		apiAddress: apiAddress,
		httpClient: http.Client{Timeout: timeout},
	}
}

func (c *HotelBedsClient) getHashSum() string {
	h := sha256.New()
	h.Write([]byte(c.apiKey + c.secret + strconv.Itoa(int(time.Now().Unix()))))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Fetches the hotels' rates for the specified bookings
func (c *HotelBedsClient) GetHotelRates(r GetHotelsRequest) (getHotels *GetHotelsResponse, request []byte, response []byte, err error) {
	request, err = json.Marshal(r)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", c.apiAddress+"/hotel-api/1.0/hotels", bytes.NewReader(request))
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Api-key", c.apiKey)
	req.Header.Add("X-Signature", c.getHashSum())

	res, err := c.httpClient.Do(req)
	if err != nil {
		return
	}

	response, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.Status != "200 OK" {
		err = errors.New("supplier error")
		return
	}

	var hotels GetHotelsResponse
	err = json.Unmarshal(response, &hotels)
	if err != nil {
		return
	}
	getHotels = &hotels
	return
}
