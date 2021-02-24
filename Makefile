ISOLATION_ID ?= local

.PHONY: build
build:
	docker build -t sxtctl:$(ISOLATION_ID) .
