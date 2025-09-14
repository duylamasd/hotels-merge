CREATE TABLE IF NOT EXISTS hotels (
  id SERIAL PRIMARY KEY,
  hotel_id TEXT UNIQUE NOT NULL,
  destination_id TEXT NOT NULL,
  name TEXT NOT NULL,
  location JSONB,
  description TEXT,
  images JSONB,
  amenities JSONB,
  booking_conditions TEXT[],
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_hotels_destination_id ON hotels(destination_id);
