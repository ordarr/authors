help:
	@echo ''
	@echo 'Usage: make [TARGET] [EXTRA_ARGUMENTS]'
	@echo 'Targets:'
	@echo 'dev: make dev for development work'
	@echo 'proto: generate gRPC definitions'

dev:
	docker compose -f docker-compose-dev.yml down
	docker compose -f docker-compose-dev.yml up --build

proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative v1/authors.proto
