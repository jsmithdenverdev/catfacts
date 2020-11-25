.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/subscription cmd/lambda/subscription/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/sendfact world/lambda/sendfact/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

