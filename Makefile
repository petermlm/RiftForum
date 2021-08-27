GOVERSION=1.17
GOBIN=go$(GOVERSION)
BIN=tmp/riftforum

migrations:
	$(GOBIN) run src/*.go -migrate true
.PHONY: migrations

run:
	$(GOBIN) run src/*.go
.PHONY: run

air:
	air -c air.toml
.PHONY: run

build:
	$(GOBIN) build -o $(BIN) src/*.go
.PHONY: build

fmt:
	$(GOBIN) fmt ./...
.PHONY: fmt
