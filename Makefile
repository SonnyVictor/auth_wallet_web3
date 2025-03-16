
DB_URL=postgresql://root:secret@localhost:5432/auth-wallet?sslmode=disable

network:
	docker network create auth-wallet

postgres:
	docker run --name auth-wallet -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it auth-wallet createdb --username=root --owner=root auth-wallet

dropdb:
	docker exec -it auth-wallet dropdb auth-wallet

migrateup:
	migrate -path db/migration -database "$(DB_URL)"  -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)"  -verbose down

migrateup1:
	migrate -path db/migration -database "$(DB_URL)"  -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_URL)"  -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server migrateup1 migratedown2 db_docs db_schema