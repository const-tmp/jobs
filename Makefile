db:
	docker run --rm -p 5432:5432 -e POSTGRES_PASSWORD=jobs postgres:latest

proto:
	protoc --proto_path=proto \
		--go_out=back/proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=back/proto \
		--go-grpc_opt=paths=source_relative proto/*.proto
	venv/bin/python -m grpc_tools.protoc \
		--proto_path=proto \
		--python_out=bot/jbot/proto \
		--grpc_python_out=bot/jbot/proto \
		proto/*.proto

up:
	docker compose up

down:
	docker compose down #--volumes

build:
	docker compose build

.PHONY: proto