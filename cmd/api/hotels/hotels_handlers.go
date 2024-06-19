package hotels

import "skabillium/liteapi/cmd/clients"

// All /hotels handlers need to share some state like the HotelBedsClient instance
// so they are all aggreagated under the HotelHanlders struct
type HotelsHandlers struct {
	hbc *clients.HotelBedsClient
}

// Instantiate a HotelHanlders struct with a HotelBeds client
func NewHotelsHanlders(hbc *clients.HotelBedsClient) *HotelsHandlers {
	return &HotelsHandlers{hbc: hbc}
}
