FROM golang:1.23-alpine AS build

ENV CGO_ENABLED=1
RUN apk add --no-cache \
    gcc \
    musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o auto cmd/auto.go

# Main Stage
FROM alpine

WORKDIR /app

COPY --from=build /app/auto .
COPY --from=build /app/config/config.yml config/config.yml

ENTRYPOINT ["./auto"]
