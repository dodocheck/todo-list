PROTO_DIR := proto
PROTO_INPUT := $(PROTO_DIR)/tasks.proto $(PROTO_DIR)/service.proto
PROTO_OUT_DIR := pkg/pb

proto-gen:
	protoc --proto_path=$(PROTO_DIR) $(PROTO_INPUT) \
		--go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative
		
deploy:
	docker compose -f deployment/docker-compose.yml build --no-cache api-service db-service logger-service
	docker compose -f deployment/docker-compose.yml up -d --force-recreate api-service db-service logger-service

down:
	docker compose -f deployment/docker-compose.yml down -v

test:
	go test ./services/api/... ./services/db/... ./services/logger/... -cover

lint:
	golangci-lint run ./services/api/... ./services/db/... ./services/logger/...

format:
	golangci-lint fmt
