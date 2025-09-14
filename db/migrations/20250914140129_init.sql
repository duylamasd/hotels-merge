-- Create "hotels" table
CREATE TABLE "hotels" (
  "id" serial NOT NULL,
  "hotel_id" text NOT NULL,
  "destination_id" text NOT NULL,
  "name" text NOT NULL,
  "location" jsonb NULL,
  "description" text NULL,
  "images" jsonb NULL,
  "amenities" jsonb NULL,
  "booking_conditions" text[] NULL,
  "created_at" timestamptz NULL DEFAULT now(),
  "updated_at" timestamptz NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "hotels_hotel_id_key" UNIQUE ("hotel_id")
);
-- Create index "idx_hotels_destination_id" to table: "hotels"
CREATE INDEX "idx_hotels_destination_id" ON "hotels" ("destination_id");
