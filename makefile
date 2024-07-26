run:
	go run ./cmd/app

build-run:
	go build -C ./cmd/app/ -o ../../recipe-api.exe && ./recipe-api.exe

live:
	nodemon --exec go run ./cmd/app --ext go

test:
	go test ./...