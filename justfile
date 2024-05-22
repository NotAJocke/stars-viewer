default: run

build:
  go build ./cmd/app/main.go

run:
  go run ./cmd/app/main.go

temp:
  go run ./cmd/temp/main.go

watch:
  air

views:
  templ generate -path ./views/

viewsw:
  templ generate -path ./views/ -watch

migrate-up:
  migrate -path migrations -database sqlite3://./db/database.db up

migrate-down:
  migrate -path migrations -database sqlite3://./db/database.db down
