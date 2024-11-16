export GOCACHE=$(SRCPATH)/tmp/go-cache

default: build

build:
	go build -o ./tmp/auto cmd/auto.go 

docker: 
	docker compose up -d
	docker compose logs -f

deps:
	go mod download
	go mod tidy

clean:
	rm -rf ./tmp/

