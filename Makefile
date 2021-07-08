MAKEFILE_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
include $(MAKEFILE_DIR)/standard_defs.mk

clean: clean_build_go

distclean: clean_docker

build: $(MARKERS)/build_go $(MARKERS)/build_docker

test: $(MARKERS)/test_go

analyze: analyze_go analyze_fossa

publish: gh-create-draft-release
	if [ "$(RELEASABLE)" = "yes" ]; then \
		$(GH_RELEASE) upload $(VERSION) target/* ; \
	fi

.PHONY: clean_docker
clean_docker: clean_toolchain_docker
	docker rmi -f sxtctl:$(ISOLATION_ID)

$(MARKERS)/build_docker: $(MARKERS)/build_go
	docker build -t sxtctl:$(ISOLATION_ID) .
	touch $@

.PHONY: analyze_dive
analyze_dive:
	$(DIVE_ANALYZE) sxtctl:$(ISOLATION_ID)
