include .env

add_new_migration:
	@atlas migrate new $(name) \
	--dir "file://db/migrations"

generate_migration:
	@atlas migrate diff $(name) \
	--dir "file://db/migrations" \
	--to "file://db/schema.sql" \
	--dev-url "docker://postgres?search_path=public"

apply_migration:
	@atlas migrate apply \
	--url $(DB_URI) \
	--dir "file://db/migrations"

revert_migration:
	@atlas migrate down \
	--dev-url "docker://postgres?search_path=public" \
	--url $(DB_URI) \
	--dir "file://db/migrations"

sqlc_generate:
	sqlc generate

run_unit_tests:
	@go test -v $$(go list ./... | grep -v ./tests/e2e)

run_e2e_tests:
	@go test -v ./tests/e2e
