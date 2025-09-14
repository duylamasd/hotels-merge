package v1

type FindHotelsQueryDTO struct {
	DestinationID *string   `form:"destination_id" binding:"omitnil,min=1"`
	HotelIDs      *[]string `form:"hotel_ids" binding:"omitnil,min=1,dive,required"`
}
