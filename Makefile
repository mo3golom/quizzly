.PHONY: build-lambda clean-lambda lint

build-lambda:
	mkdir -p functions
	GOOS=linux GOARCH=amd64 go build -o ./functions/main ./cmd/lambda/main.go
	zip -j ./functions/main.zip ./functions/main

clean-lambda:
	rm -rf functions

lint:
	golangci-lint -c ./.golangci.yml run ./...
