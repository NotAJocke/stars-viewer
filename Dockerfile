# syntax=docker/dockerfile:1.7-labs

FROM golang:1.22 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY --parents=true ./cmd/ ./internal/ ./

RUN CGO_ENABLED=1 GOOS=linux go build -o /stars-viewer ./cmd/app/main.go

#########

FROM debian:12

WORKDIR /app
COPY --parents=true ./internal/public/ ./internal/templates/ ./
COPY --from=build /stars-viewer ./stars-viewer

EXPOSE 8080
ENV CONTAINER=true

ENTRYPOINT [ "./stars-viewer" ]
