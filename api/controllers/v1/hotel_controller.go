package v1

import (
	"net/http"

	apiDomains "github.com/duylamasd/hotels-merge/api/domains"
	v1Dto "github.com/duylamasd/hotels-merge/api/dto/v1"
	"github.com/duylamasd/hotels-merge/domains"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HotelController struct {
	logger  *zap.Logger
	service domains.HotelService
}

func NewHotelController(
	logger *zap.Logger,
	service domains.HotelService,
) *HotelController {
	return &HotelController{
		logger:  logger,
		service: service,
	}
}

func (c *HotelController) Find(ctx *gin.Context) {
	var query v1Dto.FindHotelsQueryDTO
	c.logger.Info("GET /api/v1/hotels - Validating query params")
	if err := ctx.ShouldBindQuery(&query); err != nil {
		c.logger.Error(err.Error())
		e := apiDomains.NewHttpError(http.StatusBadRequest, err.Error())
		_ = ctx.Error(e)
		return
	}

	c.logger.Info("GET /api/v1/hotels - Validating either destination id or hotel ids is available")
	if query.DestinationID == nil && query.HotelIDs == nil {
		c.logger.Error("Either destination or hotel ids was not provided")
		e := apiDomains.NewHttpError(http.StatusBadRequest, "Either destination or list of hotel ids need to be provided")
		_ = ctx.Error(e)
		return
	}

	if query.DestinationID != nil {
		c.logger.Info("GET /api/v1/hotels - Finding hotels by destination id", zap.String("destination_id", *query.DestinationID))
		hotels, err := c.service.FindByDestinationID(ctx, *query.DestinationID)
		if err != nil {
			c.logger.Error("Could not fetch list of hotels by destination due to connectivity issue", zap.String("destination_id", *query.DestinationID))
			e := apiDomains.NewHttpError(http.StatusInternalServerError, "Could not fetch list of hotels. Please retry again")
			_ = ctx.Error(e)
			return
		}

		ctx.JSON(http.StatusOK, hotels)
		return
	}

	c.logger.Info("GET /api/v1/hotels - Finding hotels by hotel ids", zap.Strings("hotel_ids", *query.HotelIDs))
	hotels, err := c.service.FindByHotelIDs(ctx, *query.HotelIDs)
	if err != nil {
		c.logger.Error("Could not fetch list of hotels by list of hotel ids due to connectivity issue", zap.Strings("hotel_ids", *query.HotelIDs))
		e := apiDomains.NewHttpError(http.StatusInternalServerError, "Could not fetch list of hotels. Please retry again")
		_ = ctx.Error(e)
		return
	}

	ctx.JSON(http.StatusOK, hotels)
}
