postgres:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=secret -d postgres:17-alpine
  
createdb:
	docker exec -it postgres17 createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres17 dropdb --username=postgres simple_bank 

migrateup:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc