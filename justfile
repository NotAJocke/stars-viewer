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
  templ generate -path ./internal/templates/

viewsw:
  templ generate -path ./internal/templates/ -watch

css:
  npx tailwindcss -i ./internal/public/css/input.css -o ./internal/public/css/output.css

cssw:
  npx tailwindcss -i ./internal/public/css/input.css -o ./internal/public/css/output.css --watch

migrate-up:
  migrate -path migrations -database sqlite3://./db/database.db up

migrate-down:
  migrate -path migrations -database sqlite3://./db/database.db down
