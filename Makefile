MAKEFILE_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
include $(MAKEFILE_DIR)/standard_defs.mk

clean: clean_go

distclean: clean_docker

build: $(MARKERS)/build_docker

analyze: analyze_fossa

.PHONY: build_dev
build_dev:
	(cd ./cmd/sxtctl && go build .)

$(MARKERS)/build_go: $(MARKERS)/build_toolchain_docker
	$(TOOL) -w /project/cmd/sxtctl $(TOOLCHAIN_IMAGE) go build

.PHONY: clean_go
clean_go: $(MARKERS)/build_toolchain_docker
	$(TOOL) $(TOOLCHAIN_IMAGE) go clean

$(MARKERS)/build_docker:
	docker build -t sxtctl:$(ISOLATION_ID) .
	touch $@

.PHONY: clean_docker
clean_docker: clean_toolchain_docker
	docker rmi -f sxtctl:$(ISOLATION_ID)
