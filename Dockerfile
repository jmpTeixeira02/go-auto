FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor vendor/

COPY . .

RUN --mount=type=cache,target="/root/.cache/go-build" go build cmd/auto.go

CMD ["./auto"]
