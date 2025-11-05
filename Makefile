APP_NAME=app
DB_URL=postgresql://root:secret@localhost:5432/$(APP_NAME)?sslmode=disable

install_mac:
	npm install -g dbdocs
	npm install -g @dbml/cli
	brew install golang-migrate
	brew install sqlc
	brew install protobuf
	brew install grpc
	go install github.com/golang/mock/mockgen@latest
	go install github.com/rakyll/statik@latest
	go install google.golang.org/grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go mod tidy

test_migration: dropdb createdb migrateup migratedown

network:
	docker network create $(APP_NAME)-network

postgres:
	docker run --name postgres --network $(APP_NAME)-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine

mysql:
	docker run --name mysql8 -p 3306:3306  -e MYSQL_ROOT_PASSWORD=secret -d mysql:8

local_up:
	docker compose -f docker-compose-local.yaml up -d
	air

createdb:
	docker exec -it postgres createdb --username=root --owner=root $(APP_NAME)

dropdb:
	docker exec -it postgres dropdb $(APP_NAME) || true

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_docs:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml \
	&& dbdocs build doc/db.dbml \
	&& dbdocs password --set $(password) --project $(APP_NAME)

sqlc:
	sqlc generate

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/aria/goapp/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/aria/goapp/worker TaskDistributor

test:
	go test -v -cover -short ./...

server:
	swag init
	go run main.go

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=$(APP_NAME) \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration db_docs db_schema sqlc test server mock proto evans redis install_mac install_go_tools
