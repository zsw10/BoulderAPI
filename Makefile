BOULDER_DIR := ./boulder

.PHONY: start-boulder
start-boulder:
	@echo "Starting Boulder server..."
	@cd $(BOULDER_DIR) && docker-compose up
	@echo "Boulder server started."

.PHONY: stop-boulder
stop-boulder:
	@echo "Stopping Boulder server..."
	@cd $(BOULDER_DIR) && docker-compose down
	@echo "BOULDER server stopped."

.PHONY: start-api
start-api: 
	go run ./cmd/api

.PHONY: stop-api
stop-api:
	pkill -SIGTERM api
