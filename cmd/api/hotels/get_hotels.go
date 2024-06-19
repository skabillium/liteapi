package hotels

import (
	"encoding/json"
	"errors"
	"net/http"
	"skabillium/liteapi/cmd/clients"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// @Tags Hotels
// @Description Get rate information for specified hotels
// @Produce json
// @Success 200
// @Param checkin query string true "Check-in date"
// @Param checkout query string true "Check-out date"
// @Param currency query string true "Currency for payment"
// @Param guestNationality query string true "Nationality of guests"
// @Param hotelIds query string true "Comma separated list of hotel ids to check"
// @Param occupancies query string true "Occupancies to check, how many rooms, adults or children"
// @Router /hotels [get]
func (h *HotelsHandlers) GetHotelsHandler(c *gin.Context) {
	query, err := parseGetHotelsQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while parsing query params",
			"error":   err.Error(),
		})
		return
	}

	err = validateGetHotelsQuery(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while validating query params",
			"error":   err.Error(),
		})
		return
	}

	getHotelsReq := clients.GetHotelsRequest{
		Stay: clients.StayOpts{
			CheckIn:  query.FormatCheckin(),
			CheckOut: query.FormatCheckout(),
		},
		Occupancies: query.Occupancies,
		Hotels: clients.HotelOpts{
			Hotel: query.HotelIds,
		},
	}

	getHotels, request, response, err := h.hbc.GetHotelRates(getHotelsReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while connecting to supplier API",
			"error":   err.Error(),
			"supplier": gin.H{
				"request":  string(request),
				"response": string(response),
			},
		})
		return
	}

	// Find the cheapest for every hotel
	data := make([]HotelData, getHotels.Hotels.Total)
	for i := 0; i < len(getHotels.Hotels.Hotels); i++ {
		hotel := getHotels.Hotels.Hotels[i]
		price, err := strconv.ParseFloat(hotel.MinRate, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error while parsing to supplier API response",
				"error":   err.Error(),
				"supplier": gin.H{
					"request":  string(request),
					"response": string(response),
				},
			})
			return
		}

		data[i].Currency = hotel.Currency
		data[i].HotelId = strconv.Itoa(hotel.Code)
		data[i].Price = float32(price)
	}

	// Format the response according to spec
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"supplier": gin.H{
			"request":  string(request),
			"response": string(response),
		},
	})
}

// Struct that holds the response data
type HotelData struct {
	HotelId  string  `json:"hotelId"`
	Currency string  `json:"currency"`
	Price    float32 `json:"price"`
}

// Allowed query params
type GetHotelsQuery struct {
	CheckIn          time.Time
	CheckOut         time.Time
	HotelIds         []int
	Occupancies      []clients.Occupancy
	Currency         string
	GuestNationality string
}

func (q *GetHotelsQuery) FormatCheckin() string {
	return q.CheckIn.Format(time.DateOnly)
}

func (q *GetHotelsQuery) FormatCheckout() string {
	return q.CheckOut.Format(time.DateOnly)
}

// Read request query params into a GetHotelsQuery instance
func parseGetHotelsQuery(c *gin.Context) (*GetHotelsQuery, error) {
	q := &GetHotelsQuery{}
	checkin, err := time.Parse(time.DateOnly, c.Query("checkin"))
	if err != nil {
		return nil, err
	}
	checkout, err := time.Parse(time.DateOnly, c.Query("checkout"))
	if err != nil {
		return nil, err
	}
	q.CheckIn = checkin
	q.CheckOut = checkout

	// Parse hotel ids as integers
	hotelsStr := strings.Split(c.Query("hotelIds"), ",")
	q.HotelIds = make([]int, len(hotelsStr))
	for i := 0; i < len(hotelsStr); i++ {
		id, err := strconv.Atoi(hotelsStr[i])
		if err != nil {
			return nil, err
		}
		q.HotelIds[i] = id
	}

	// Parse occupancies json
	err = json.Unmarshal([]byte(c.Query("occupancies")), &q.Occupancies)
	if err != nil {
		return nil, err
	}
	return q, nil
}

// Performs validation on a GetHotelsQuery
func validateGetHotelsQuery(q *GetHotelsQuery) error {
	now := time.Now()
	if now.After(q.CheckIn) {
		return errors.New("checkin must be in the future")
	}
	if q.CheckIn.After(q.CheckOut) {
		return errors.New("checkin must be after checkout")
	}

	if len(q.HotelIds) == 0 {
		return errors.New("at least one hotel id is required")
	}

	for i := 0; i < len(q.Occupancies); i++ {
		if q.Occupancies[i].Rooms < 1 {
			return errors.New("rooms must be a positive integer")
		}
		if q.Occupancies[i].Adults < 0 {
			return errors.New("adults must be a non negative integer")
		}
		if q.Occupancies[i].Children < 0 {
			return errors.New("children must be a non negative integer")
		}
	}

	return nil
}
