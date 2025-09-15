package v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "github.com/duylamasd/hotels-merge/api/controllers/v1"
	"github.com/duylamasd/hotels-merge/api/domains"
	"github.com/duylamasd/hotels-merge/api/middlewares"
	"github.com/duylamasd/hotels-merge/config"
	"github.com/duylamasd/hotels-merge/lib"
	"github.com/duylamasd/hotels-merge/mocks"
	"github.com/duylamasd/hotels-merge/sqlc"
	"github.com/duylamasd/hotels-merge/sqlc/dto"
	"github.com/gin-gonic/gin"
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

func TestHotelController_Find(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, _ := lib.NewLogger(&config.Config{LogLevel: "info"})
	mockHotelService := mocks.NewMockHotelService(ctrl)
	hotelController := v1.NewHotelController(logger, mockHotelService)
	errorHandler := middlewares.NewErrorHandler(logger)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(errorHandler.Handler())

	api := router.Group("/api/v1")
	hotels := api.Group("/hotels")
	hotels.GET("", hotelController.Find)

	t.Run("should return 400 if neither destination_id nor hotel_ids is provided", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/hotels", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response domains.HttpError
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, "Either destination or list of hotel ids need to be provided", response.Message)
	})

	t.Run("should return 200 with list of hotels when destination_id is provided", func(t *testing.T) {
		destinationID := "dest_456"

		expectedHotels := []*sqlc.Hotel{
			{ID: 1, HotelID: "hotel_123", DestinationID: destinationID, Name: "Test Hotel 1", Location: createMockLocation(), Description: nil, Images: nil, Amenities: nil, BookingConditions: []string{"No smoking"}, CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}},
			{ID: 2, HotelID: "hotel_124", DestinationID: destinationID, Name: "Test Hotel 2", Location: createMockLocation(), Description: nil, Images: nil, Amenities: nil, BookingConditions: []string{"No pets"}, CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}},
		}

		mockHotelService.EXPECT().FindByDestinationID(gomock.Any(), destinationID).Return(expectedHotels, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/hotels?destination_id="+destinationID, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*sqlc.Hotel
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Len(t, response, len(expectedHotels))
		assert.Equal(t, response[0].DestinationID, destinationID)
		assert.Equal(t, response[1].DestinationID, destinationID)
	})

	t.Run("should return 200 with list of hotels when hotel_ids is provided", func(t *testing.T) {
		hotelIDs := []string{"hotel_123", "hotel_124"}
		destinationID := "dest_456"

		expectedHotels := []*sqlc.Hotel{
			{ID: 1, HotelID: "hotel_123", DestinationID: destinationID, Name: "Test Hotel 1", Location: createMockLocation(), Description: nil, Images: nil, Amenities: nil, BookingConditions: []string{"No smoking"}, CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}},
			{ID: 2, HotelID: "hotel_124", DestinationID: destinationID, Name: "Test Hotel 2", Location: createMockLocation(), Description: nil, Images: nil, Amenities: nil, BookingConditions: []string{"No pets"}, CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}},
		}

		mockHotelService.EXPECT().FindByHotelIDs(gomock.Any(), hotelIDs).Return(expectedHotels, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/hotels?hotel_ids=hotel_123&hotel_ids=hotel_124", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*sqlc.Hotel
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Len(t, response, len(expectedHotels))
		assert.Equal(t, response[0].HotelID, expectedHotels[0].HotelID)
		assert.Equal(t, response[1].HotelID, expectedHotels[1].HotelID)
	})

	t.Run("should return 400 if destination_id is empty string", func(t *testing.T) {
		destinationID := ""

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/hotels?destination_id="+destinationID, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response domains.HttpError
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, "Key: 'FindHotelsQueryDTO.DestinationID' Error:Field validation for 'DestinationID' failed on the 'min' tag", response.Message)
	})

	t.Run("should return 400 if hotel_ids is empty", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/hotels?hotel_ids=", nil)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response domains.HttpError
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("should return 400 if any of hotel_ids is empty string", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/hotels?hotel_ids=hotel_123&hotel_ids=", nil)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response domains.HttpError
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("should return 500 if service returns error for destination_id", func(t *testing.T) {
		destinationID := "dest_error"

		mockHotelService.EXPECT().FindByDestinationID(gomock.Any(), destinationID).Return(nil, pgx.ErrTxClosed).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/hotels?destination_id="+destinationID, nil)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response domains.HttpError
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, "Could not fetch list of hotels. Please retry again", response.Message)
	})

	t.Run("should return 500 if service returns error for hotel_ids", func(t *testing.T) {
		mockHotelService.EXPECT().FindByHotelIDs(gomock.Any(), []string{"hotel_error"}).Return(nil, pgx.ErrTxClosed).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/hotels?hotel_ids=hotel_error", nil)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response domains.HttpError
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, "Could not fetch list of hotels. Please retry again", response.Message)
	})

	t.Run("should return 200 with empty list when no hotels found for destination_id", func(t *testing.T) {
		destinationID := "dest_no_hotels"

		mockHotelService.EXPECT().FindByDestinationID(gomock.Any(), destinationID).Return([]*sqlc.Hotel{}, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/hotels?destination_id="+destinationID, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*sqlc.Hotel
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Len(t, response, 0)
	})

	t.Run("should return 200 with empty list when no hotels found for hotel_ids", func(t *testing.T) {
		mockHotelService.EXPECT().FindByHotelIDs(gomock.Any(), []string{"non_existent_hotel"}).Return([]*sqlc.Hotel{}, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/hotels?hotel_ids=non_existent_hotel", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*sqlc.Hotel
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Len(t, response, 0)
	})

	t.Run("should return 200 with list of existing hotels when some hotel_ids do not exist", func(t *testing.T) {
		destinationID := "dest_456"
		expectedHotels := []*sqlc.Hotel{
			{ID: 1, HotelID: "hotel_123", DestinationID: destinationID, Name: "Test Hotel 1", Location: createMockLocation(), Description: nil, Images: nil, Amenities: nil, BookingConditions: []string{"No smoking"}, CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}},
		}

		mockHotelService.EXPECT().FindByHotelIDs(gomock.Any(), []string{"hotel_123", "non_existent_hotel"}).Return(expectedHotels, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/hotels?hotel_ids=hotel_123&hotel_ids=non_existent_hotel", nil)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var response []*sqlc.Hotel
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Len(t, response, len(expectedHotels))
		assert.Equal(t, response[0].HotelID, expectedHotels[0].HotelID)
	})

	t.Run("should return 200 of list of hotels by destination_id when both destination_id and hotel_ids are provided", func(t *testing.T) {
		destinationID := "dest_456"

		expectedHotels := []*sqlc.Hotel{
			{ID: 1, HotelID: "hotel_123", DestinationID: destinationID, Name: "Test Hotel 1", Location: createMockLocation(), Description: nil, Images: nil, Amenities: nil, BookingConditions: []string{"No smoking"}, CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}},
			{ID: 2, HotelID: "hotel_124", DestinationID: destinationID, Name: "Test Hotel 2", Location: createMockLocation(), Description: nil, Images: nil, Amenities: nil, BookingConditions: []string{"No pets"}, CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}},
		}

		mockHotelService.EXPECT().FindByDestinationID(gomock.Any(), destinationID).Return(expectedHotels, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/hotels?destination_id="+destinationID+"&hotel_ids=hotel_123&hotel_ids=hotel_124", nil)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var response []*sqlc.Hotel
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Len(t, response, len(expectedHotels))
		assert.Equal(t, response[0].DestinationID, destinationID)
		assert.Equal(t, response[1].DestinationID, destinationID)
	})
}
