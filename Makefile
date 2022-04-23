GO ?= go
GOTEST ?= $(GO) test

.PHONY: test generate/mocks

test:
	$(GOTEST) -p 1 ./...


generate/mocks:
	go generate ./...