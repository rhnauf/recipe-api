run:
	go run ./cmd/app

build-run:
	go build -C ./cmd/app/ -o ../../recipe-api && ./recipe-api

live:
	nodemon --exec go run ./cmd/app --ext go

test:
	go test ./...
