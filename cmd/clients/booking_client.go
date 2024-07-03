package clients

// Generic client interface to help with testing
type BookingClient interface {
	GetHotelRates(r GetHotelsRequest) (*GetHotelsResponse, []byte, []byte, error)
}

type StayOpts struct {
	CheckIn  string `json:"checkIn"`
	CheckOut string `json:"checkOut"`
}

type HotelOpts struct {
	Hotel []int `json:"hotel"`
}

type Occupancy struct {
	Rooms    int `json:"rooms"`
	Adults   int `json:"adults"`
	Children int `json:"children"`
}

// Request structure for GET /hotels
type GetHotelsRequest struct {
	Stay        StayOpts    `json:"stay"`
	Occupancies []Occupancy `json:"occupancies"`
	Hotels      HotelOpts   `json:"hotels"`
}

// Response structure for GET /hotels
type GetHotelsResponse struct {
	Hotels struct {
		Total  int `json:"total"`
		Hotels []struct {
			Code     int    `json:"code"`
			Currency string `json:"currency"`
			MinRate  string `json:"minRate"`
		} `json:"hotels"`
	} `json:"hotels"`
}
