.PHONY: run migrate seed

run:
	go run cmd/threter-ticket-system/main.go

migrate:
	@echo "Running migrations..."
	@go run cmd/threter-ticket-system/scripts/migrate/main.go

seed:
	@echo "Seeding database..."
	@go run cmd/threter-ticket-system/scripts/seed/main.go
