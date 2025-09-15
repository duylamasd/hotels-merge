package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/duylamasd/hotels-merge/config"
	"github.com/duylamasd/hotels-merge/lib"
	"github.com/duylamasd/hotels-merge/mocks"
	"github.com/duylamasd/hotels-merge/services"
	"github.com/duylamasd/hotels-merge/sqlc"
	"github.com/duylamasd/hotels-merge/sqlc/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func createMockLocation() *dto.HotelLocation {
	latitude := 10.762622
	longitude := 106.660172
	address := "123 Test St, Test City"
	city := "Test City"
	country := "Test Country"

	return &dto.HotelLocation{
		Latitude:  &latitude,
		Longitude: &longitude,
		Address:   &address,
		City:      &city,
		Country:   &country,
	}
}

func TestHotelService_FindByHotelID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, _ := lib.NewLogger(&config.Config{LogLevel: "info"})
	mockSqlcQuerier := mocks.NewMockQuerier(ctrl)
	hotelService := services.NewHotelService(logger, &config.DBStore{
		Queries:  mockSqlcQuerier,
		ConnPool: nil,
	})

	t.Run("should return hotel when found", func(t *testing.T) {
		ctx := context.Background()

		hotelID := "hotel_123"

		expectedHotel := &sqlc.Hotel{
			ID:                1,
			HotelID:           hotelID,
			DestinationID:     "dest_456",
			Name:              "Test Hotel",
			Location:          createMockLocation(),
			Description:       nil,
			Images:            nil,
			Amenities:         nil,
			BookingConditions: []string{"No smoking", "No pets"},
			CreatedAt:         pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt:         pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		mockSqlcQuerier.
			EXPECT().
			FindHotelByHotelID(ctx, hotelID).
			Return(expectedHotel, nil).
			Times(1)

		result, err := hotelService.FindByHotelID(ctx, hotelID)

		assert.NoError(t, err)
		assert.Equal(t, result.ID, expectedHotel.ID)
		assert.Equal(t, result.HotelID, expectedHotel.HotelID)
		assert.Equal(t, result.DestinationID, expectedHotel.DestinationID)
		assert.Equal(t, result.Name, expectedHotel.Name)
	})

	t.Run("should return error when hotel not found", func(t *testing.T) {
		ctx := context.Background()

		hotelID := "non_existent_hotel"

		mockSqlcQuerier.EXPECT().FindHotelByHotelID(ctx, hotelID).Return(nil, pgx.ErrNoRows).Times(1)

		result, err := hotelService.FindByHotelID(ctx, hotelID)

		assert.Error(t, err)
		assert.Equal(t, err, pgx.ErrNoRows)
		assert.Nil(t, result)
	})
}

func TestHotelService_FindByDestinationID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, _ := lib.NewLogger(&config.Config{LogLevel: "info"})
	mockSqlcQuerier := mocks.NewMockQuerier(ctrl)
	hotelService := services.NewHotelService(logger, &config.DBStore{
		Queries:  mockSqlcQuerier,
		ConnPool: nil,
	})

	t.Run("should return hotels when found", func(t *testing.T) {
		ctx := context.Background()
		destinationID := "dest_456"

		expectedHotels := []*sqlc.Hotel{
			{ID: 1, HotelID: "hotel_123", DestinationID: destinationID, Name: "Test Hotel 1", Location: createMockLocation(), Description: nil, Images: nil, Amenities: nil, BookingConditions: []string{"No smoking"}, CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}},
			{ID: 2, HotelID: "hotel_124", DestinationID: destinationID, Name: "Test Hotel 2", Location: createMockLocation(), Description: nil, Images: nil, Amenities: nil, BookingConditions: []string{"No pets"}, CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}},
		}

		mockSqlcQuerier.EXPECT().FindHotelsByDestinationID(ctx, destinationID).Return(expectedHotels, nil).Times(1)

		result, err := hotelService.FindByDestinationID(ctx, destinationID)

		assert.NoError(t, err)
		assert.Len(t, result, len(expectedHotels))
		assert.Equal(t, result[0].DestinationID, destinationID)
		assert.Equal(t, result[1].DestinationID, destinationID)
	})

	t.Run("should return empty list when no hotels found", func(t *testing.T) {
		ctx := context.Background()
		destinationID := "empty_dest"

		expectedHotels := []*sqlc.Hotel{}

		mockSqlcQuerier.EXPECT().FindHotelsByDestinationID(ctx, destinationID).Return(expectedHotels, nil).Times(1)

		result, err := hotelService.FindByDestinationID(ctx, destinationID)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("should return error when query fails", func(t *testing.T) {
		ctx := context.Background()
		destinationID := "error_dest"

		mockSqlcQuerier.EXPECT().FindHotelsByDestinationID(ctx, destinationID).Return(nil, pgx.ErrTxClosed).Times(1)

		result, err := hotelService.FindByDestinationID(ctx, destinationID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestHotelService_FindByHotelIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, _ := lib.NewLogger(&config.Config{LogLevel: "info"})
	mockSqlcQuerier := mocks.NewMockQuerier(ctrl)
	hotelService := services.NewHotelService(logger, &config.DBStore{
		Queries:  mockSqlcQuerier,
		ConnPool: nil,
	})

	t.Run("should return hotels when found", func(t *testing.T) {
		ctx := context.Background()
		hotelIDs := []string{"hotel_123", "hotel_124"}
		destinationID := "dest_456"

		expectedHotels := []*sqlc.Hotel{
			{ID: 1, HotelID: "hotel_123", DestinationID: destinationID, Name: "Test Hotel 1", Location: createMockLocation(), Description: nil, Images: nil, Amenities: nil, BookingConditions: []string{"No smoking"}, CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}},
			{ID: 2, HotelID: "hotel_124", DestinationID: destinationID, Name: "Test Hotel 2", Location: createMockLocation(), Description: nil, Images: nil, Amenities: nil, BookingConditions: []string{"No pets"}, CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}},
		}

		mockSqlcQuerier.EXPECT().FindHotelsByHotelIDs(ctx, hotelIDs).Return(expectedHotels, nil).Times(1)

		result, err := hotelService.FindByHotelIDs(ctx, hotelIDs)

		assert.NoError(t, err)
		assert.Len(t, result, len(expectedHotels))
		assert.Equal(t, result[0].HotelID, expectedHotels[0].HotelID)
		assert.Equal(t, result[1].HotelID, expectedHotels[1].HotelID)
	})

	t.Run("should return empty list when no hotels found", func(t *testing.T) {
		ctx := context.Background()
		hotelIDs := []string{"non_existent_hotel"}

		mockSqlcQuerier.EXPECT().FindHotelsByHotelIDs(ctx, hotelIDs).Return([]*sqlc.Hotel{}, nil).Times(1)

		result, err := hotelService.FindByHotelIDs(ctx, hotelIDs)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("should return only hotels existing in DB", func(t *testing.T) {
		ctx := context.Background()
		hotelIDs := []string{"hotel_123", "non_existent_hotel"}
		destinationID := "dest_456"

		expectedHotels := []*sqlc.Hotel{
			{ID: 1, HotelID: "hotel_123", DestinationID: destinationID, Name: "Test Hotel 1", Location: createMockLocation(), Description: nil, Images: nil, Amenities: nil, BookingConditions: []string{"No smoking"}, CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}},
		}

		mockSqlcQuerier.EXPECT().FindHotelsByHotelIDs(ctx, hotelIDs).Return(expectedHotels, nil).Times(1)

		result, err := hotelService.FindByHotelIDs(ctx, hotelIDs)

		assert.NoError(t, err)
		assert.Len(t, result, len(expectedHotels))
		assert.Equal(t, result[0].HotelID, hotelIDs[0])
		assert.NotEqual(t, result[0].HotelID, hotelIDs[1])
	})

	t.Run("should return error when query fails", func(t *testing.T) {
		ctx := context.Background()
		hotelIDs := []string{"hotel_123", "hotel_124"}

		mockSqlcQuerier.EXPECT().FindHotelsByHotelIDs(ctx, hotelIDs).Return(nil, pgx.ErrTxClosed).Times(1)

		result, err := hotelService.FindByHotelIDs(ctx, hotelIDs)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
