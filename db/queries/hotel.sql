-- name: FindHotelByHotelID :one
SELECT *
FROM hotels
WHERE hotel_id = $1;

-- name: FindHotelsByDestinationID :many
SELECT *
FROM hotels
WHERE destination_id = $1;

-- name: FindHotelsByHotelIDs :many
SELECT *
FROM hotels
WHERE hotel_id = ANY(sqlc.arg('hotel_ids')::TEXT[]);

-- name: FindHotelsByDestinationAndHotelIDs :many
SELECT *
FROM hotels
WHERE destination_id = sqlc.arg('destination_id')
  AND hotel_id = ANY(sqlc.arg('hotel_ids')::TEXT[]);
