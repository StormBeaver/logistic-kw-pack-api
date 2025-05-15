.PHONY: build
build:
	go build cmd/logistic-kw-pack-api/main.go

.PHONY: test
test:
	go test -v ./...