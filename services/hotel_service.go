package services

import (
	"context"

	"github.com/duylamasd/hotels-merge/config"
	"github.com/duylamasd/hotels-merge/domains"
	"github.com/duylamasd/hotels-merge/sqlc"
	"go.uber.org/zap"
)

type hotelService struct {
	logger *zap.Logger
	db     *config.DBStore
}

func NewHotelService(logger *zap.Logger, db *config.DBStore) domains.HotelService {
	return &hotelService{
		logger: logger,
		db:     db,
	}
}

func (s *hotelService) FindByHotelID(ctx context.Context, hotelID string) (*sqlc.Hotel, error) {
	return s.db.Queries.FindHotelByHotelID(ctx, hotelID)
}

func (s *hotelService) FindByDestinationID(ctx context.Context, destinationID string) ([]*sqlc.Hotel, error) {
	return s.db.Queries.FindHotelsByDestinationID(ctx, destinationID)
}

func (s *hotelService) FindByHotelIDs(ctx context.Context, hotelIDs []string) ([]*sqlc.Hotel, error) {
	return s.db.Queries.FindHotelsByHotelIDs(ctx, hotelIDs)
}

func (s *hotelService) FindByDestinationAndHotelIDs(ctx context.Context, destinationID string, hotelIDs []string) ([]*sqlc.Hotel, error) {
	return s.db.Queries.FindHotelsByDestinationAndHotelIDs(ctx, sqlc.FindHotelsByDestinationAndHotelIDsParams{
		DestinationID: destinationID,
		HotelIds:      hotelIDs,
	})
}
