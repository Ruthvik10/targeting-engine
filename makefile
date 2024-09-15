test:
	@go test $(shell go list ./... | grep -v 'github.com/Ruthvik10/targeting-engine/cache/mock' | grep -v 'github.com/Ruthvik10/targeting-engine/cache' | grep -v 'github.com/Ruthvik10/targeting-engine/model' | grep -v 'github.com/Ruthvik10/targeting-engine/store/mock') -count=1 -cover
run:
	@go run main.go
compose-up:
	@docker compose up --build -d
compose-down:
	@docker compose down

.PHONY: test run compose-up compose-down