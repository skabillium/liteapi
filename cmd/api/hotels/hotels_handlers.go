package hotels

import (
	"skabillium/liteapi/cmd/clients"
	"sync"

	"github.com/gin-gonic/gin"
)

// Register endpoints for "/hotels" group, also handle access control for booking client using a mutex
func RegisterHotelHandlers(r *gin.Engine, client clients.BookingClient) {
	var mu sync.RWMutex
	hotelsGroup := r.Group("/hotels")

	hotelsGroup.GET("/", func(ctx *gin.Context) { GetHotelsHandler(ctx, client, &mu) })
}
