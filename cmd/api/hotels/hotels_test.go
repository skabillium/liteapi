package hotels_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"skabillium/liteapi/cmd/api/hotels"
	"skabillium/liteapi/cmd/clients"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Encode GetHotelsQuery params to a URL
func encodeGetHotelsQuery(u *url.URL, q hotels.GetHotelsQuery) error {
	params := url.Values{}
	params.Add("checkin", q.FormatCheckin())
	params.Add("checkout", q.FormatCheckout())
	params.Add("currency", q.Currency)
	params.Add("guestNationality", q.GuestNationality)

	var hotelIds = []string{}
	for _, h := range q.HotelIds {
		hotelIds = append(hotelIds, strconv.Itoa(h))
	}
	params.Add("hotelIds", strings.Join(hotelIds, ","))

	occupanciesJson, err := json.Marshal(q.Occupancies)
	if err != nil {
		return fmt.Errorf("Could not encode occupancies %+v", q.Occupancies)
	}

	params.Add("occupancies", string(occupanciesJson))
	u.RawQuery = params.Encode()
	return nil
}

// Generic test for GET /hotels endpoint. For now only contains some invalid query params
// and expects a 400 response. Could be extended with more cases.
func TestGetHotels(t *testing.T) {
	location, err := time.LoadLocation("Local")
	if err != nil {
		t.Error("Error while loading location")
	}
	invalidQueries := []hotels.GetHotelsQuery{
		// Check-in after check-out
		{
			CheckIn:          time.Date(2024, 7, 22, 0, 0, 0, 0, location),
			CheckOut:         time.Date(2024, 7, 20, 0, 0, 0, 0, location),
			HotelIds:         []int{1, 2, 3, 4, 5},
			Currency:         "USD",
			GuestNationality: "USA",
			Occupancies: []clients.Occupancy{
				{Rooms: 2, Adults: 2},
			},
		},
		// Negative occupancy values
		{
			CheckIn:          time.Date(2024, 7, 22, 0, 0, 0, 0, location),
			CheckOut:         time.Date(2024, 7, 20, 0, 0, 0, 0, location),
			HotelIds:         []int{1, 2, 3, 4, 5},
			Currency:         "USD",
			GuestNationality: "USA",
			Occupancies: []clients.Occupancy{
				{Rooms: -2, Adults: -2},
			},
		},
	}

	urlStr := "http://localhost:8080/hotels"
	httpClient := &http.Client{}
	for _, q := range invalidQueries {
		u, err := url.Parse(urlStr)
		if err != nil {
			t.Errorf("Error while parsing url %s", urlStr)
		}

		encodeGetHotelsQuery(u, q)
		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			t.Error("Error while initializing request", err)
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			t.Error("Error while sending request", err)
		}

		status := resp.Status[:3]
		if status != "400" {
			t.Error("Received unexpected status code", status)
		}
		resp.Body.Close()
	}
}
