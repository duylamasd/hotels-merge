package v1

import (
	v1Controllers "github.com/duylamasd/hotels-merge/api/controllers/v1"
	"github.com/gin-gonic/gin"
)

type HotelRoutes struct {
	controller *v1Controllers.HotelController
}

func (s *HotelRoutes) Register(group *gin.RouterGroup) {
	hotels := group.Group("/hotels")
	hotels.GET("", s.controller.Find)
}

func NewHotelRoutes(
	controller *v1Controllers.HotelController,
) *HotelRoutes {
	return &HotelRoutes{
		controller: controller,
	}
}
