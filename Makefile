ISOLATION_ID ?= local

.PHONY: dev
dev:
	(cd ./cmd/sxtctl && go build .)

.PHONY: build
build:
	docker build -t sxtctl:$(ISOLATION_ID) .

.PHONY: all
all:

.PHONY: clean
clean:

.PHONY: test
test:
