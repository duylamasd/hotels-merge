package dto

type HotelImage struct {
	Link        string `json:"link"`
	Description string `json:"description"`
}

type HotelImages struct {
	Rooms     []HotelImage `json:"rooms"`
	Services  []HotelImage `json:"services"`
	Amenities []HotelImage `json:"amenities"`
}

type HotelAmenities struct {
	General []string `json:"general"`
	Room    []string `json:"room"`
}

type HotelLocation struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Address   *string  `json:"address"`
	City      *string  `json:"city"`
	Country   *string  `json:"country"`
}
