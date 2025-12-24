PROTO_DIR := proto
PROTO_INPUT := $(PROTO_DIR)/*.proto
PROTO_OUT_DIR := pkg/pb

protogen:
	protoc --proto_path=$(PROTO_DIR) $(PROTO_INPUT) \
		--go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
		
gen-db:
	@$(call MKDIR_P,$(DB_PB_DIR))
	protoc --proto_path=$(PROTO_DIR) $(PROTO_INPUT) \
		--go_out=$(DB_PB_DIR) --go_opt=paths=source_relative \
		--go_opt=Mtasks.proto=$(DB_GO_PKG) --go_opt=Mproto/tasks.proto=$(DB_GO_PKG) \
		--go_opt=Mservice.proto=$(DB_GO_PKG) --go_opt=Mproto/service.proto=$(DB_GO_PKG) \
		--go-grpc_out=$(DB_PB_DIR) --go-grpc_opt=paths=source_relative \
		--go-grpc_opt=Mtasks.proto=$(DB_GO_PKG) --go-grpc_opt=Mproto/tasks.proto=$(DB_GO_PKG) \
		--go-grpc_opt=Mservice.proto=$(DB_GO_PKG) --go-grpc_opt=Mproto/service.proto=$(DB_GO_PKG)

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
