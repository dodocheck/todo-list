.PHONY: gen gen-api gen-db clean clean-api clean-db

PROTO_DIR := proto
PROTO_INPUT := $(PROTO_DIR)/tasks.proto $(PROTO_DIR)/service.proto

API_PB_DIR := services/api/pb
DB_PB_DIR  := services/db/pb

API_GO_PKG := github.com/dodocheck/go-pet-project-1/services/api/pb
DB_GO_PKG  := github.com/dodocheck/go-pet-project-1/services/db/pb

# mkdir/rm portable (Windows cmd -> PowerShell)
ifeq ($(OS),Windows_NT)
  define MKDIR_P
    powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path '$(1)' | Out-Null"
  endef
  define RM_RF
    powershell -NoProfile -Command "Remove-Item -Recurse -Force -ErrorAction SilentlyContinue '$(1)'"
  endef
else
  define MKDIR_P
    mkdir -p $(1)
  endef
  define RM_RF
    rm -rf $(1)
  endef
endif

gen: gen-api gen-db

gen-api:
	@$(call MKDIR_P,$(API_PB_DIR))
	protoc --proto_path=$(PROTO_DIR) $(PROTO_INPUT) \
		--go_out=$(API_PB_DIR) --go_opt=paths=source_relative \
		--go_opt=Mtasks.proto=$(API_GO_PKG) --go_opt=Mproto/tasks.proto=$(API_GO_PKG) \
		--go_opt=Mservice.proto=$(API_GO_PKG) --go_opt=Mproto/service.proto=$(API_GO_PKG) \
		--go-grpc_out=$(API_PB_DIR) --go-grpc_opt=paths=source_relative \
		--go-grpc_opt=Mtasks.proto=$(API_GO_PKG) --go-grpc_opt=Mproto/tasks.proto=$(API_GO_PKG) \
		--go-grpc_opt=Mservice.proto=$(API_GO_PKG) --go-grpc_opt=Mproto/service.proto=$(API_GO_PKG)

gen-db:
	@$(call MKDIR_P,$(DB_PB_DIR))
	protoc --proto_path=$(PROTO_DIR) $(PROTO_INPUT) \
		--go_out=$(DB_PB_DIR) --go_opt=paths=source_relative \
		--go_opt=Mtasks.proto=$(DB_GO_PKG) --go_opt=Mproto/tasks.proto=$(DB_GO_PKG) \
		--go_opt=Mservice.proto=$(DB_GO_PKG) --go_opt=Mproto/service.proto=$(DB_GO_PKG) \
		--go-grpc_out=$(DB_PB_DIR) --go-grpc_opt=paths=source_relative \
		--go-grpc_opt=Mtasks.proto=$(DB_GO_PKG) --go-grpc_opt=Mproto/tasks.proto=$(DB_GO_PKG) \
		--go-grpc_opt=Mservice.proto=$(DB_GO_PKG) --go-grpc_opt=Mproto/service.proto=$(DB_GO_PKG)

clean: clean-api clean-db

clean-api:
	@$(call RM_RF,$(API_PB_DIR))

clean-db:
	@$(call RM_RF,$(DB_PB_DIR))

deploy:
	docker compose -f deployment/docker-compose.yml build --no-cache api-service db-service logger-service
	docker compose -f deployment/docker-compose.yml up -d --force-recreate api-service db-service logger-service

down:
	docker compose -f deployment/docker-compose.yml down -v
