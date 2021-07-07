deps:
	go get -d -t ./...

test: deps
	go test -v

bench: deps
	go test -bench . -benchmem -benchtime 5s -count 10

build: deps
	CGO_ENABLED=0 go build

lint:
	go vet
	golint -set_exit_status
