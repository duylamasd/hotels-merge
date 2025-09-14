package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type V1Routes struct {
	HotelRoutes *HotelRoutes
}

func (r *V1Routes) Register(group *gin.RouterGroup) {
	r.HotelRoutes.Register(group)
}

func NewV1Routes(
	hotelRoutes *HotelRoutes,
) *V1Routes {
	return &V1Routes{
		HotelRoutes: hotelRoutes,
	}
}

var Module = fx.Options(
	fx.Provide(NewHotelRoutes),
	fx.Provide(NewV1Routes),
)
