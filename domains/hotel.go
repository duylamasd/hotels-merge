package domains

import (
	"context"

	"github.com/duylamasd/hotels-merge/sqlc"
)

type HotelService interface {
	FindByHotelID(ctx context.Context, hotelID string) (*sqlc.Hotel, error)
	FindByDestinationID(ctx context.Context, destinationID string) ([]*sqlc.Hotel, error)
	FindByHotelIDs(ctx context.Context, hotelIDs []string) ([]*sqlc.Hotel, error)
	FindByDestinationAndHotelIDs(ctx context.Context, destinationID string, hotelIDs []string) ([]*sqlc.Hotel, error)
}
