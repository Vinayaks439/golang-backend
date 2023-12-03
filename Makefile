createdb:
	docker exec -ti postgres createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -ti postgres dropdb simple_bank
postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
migrateinit:
	migrate create -ext sql -dir db/migration -seq init_schema
migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
removepg:
	docker kill postgres && docker rm postgres
.PHONY: postgres createdb dropdb migrateup migratedown removepg

