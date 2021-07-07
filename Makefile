VERSION := $(shell git tag | grep ^v | sort -V | tail -n 1)
deps:
	go get -d -t ./...

test: deps
	go test -v

bench: deps
	go test -bench . -benchmem -benchtime 5s -count 10

build: deps
	CGO_ENABLED=0 go build -o ./bin/oyaki -ldflags "-X main.version=$(VERSION)"

lint:
	go vet
	golint -set_exit_status
