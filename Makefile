.DEFAULT_GOAL := test

.PHONY: build-example
build-example:
	go build -o example examples/main.go
	cd v4; go build -o example-v4 examples/main.go

.PHONY: run-example
run-example:
	go run examples/main.go

.PHONY: run-example-v4
run-example-v4:
	cd v4; go run examples/main.go

.PHONY: test
test:
	go test -count 1 -race -p 1 ./...
	cd v4; go test -count 1 -race -p 1 ./...
