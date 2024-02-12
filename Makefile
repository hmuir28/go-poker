build:
	@go build -o ./bin/go-poker

run: build
	@./bin/go-poker

test:
	go test -v ./...
