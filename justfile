default: run

build:
  go build ./cmd/app/main.go

run:
  go run ./cmd/app/main.go

temp:
  go run ./cmd/temp/main.go

watch:
  air
