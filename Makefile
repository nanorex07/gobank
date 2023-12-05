postgres:
	docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root gobank

listtables:
	docker exec -it postgres psql -U root -d gobank -c "\dt"

dropdb:
	docker exec -it postgres dropdb gobank

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/gobank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/gobank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb migrateup migratedown listtables sqlc