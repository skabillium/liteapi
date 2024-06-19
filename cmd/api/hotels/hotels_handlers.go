package hotels

import "skabillium/liteapi/cmd/clients"

type HotelsHandlers struct {
	hbc *clients.HotelBedsClient
}

func NewHotelsHanlders(hbc *clients.HotelBedsClient) *HotelsHandlers {
	return &HotelsHandlers{hbc: hbc}
}
