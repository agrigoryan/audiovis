.PHONY: build
build: main.go
	go build -o audiovis main.go

.PHONY: run
run: build
	./audiovis

.PHONY: test
test:
	go test ./...
