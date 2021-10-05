MAKEFILE_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
include $(MAKEFILE_DIR)/standard_defs.mk

S3_BASE_URL := s3://sxtctl

clean: clean_build_go

distclean: clean_docker

build: $(MARKERS)/build_go $(MARKERS)/build_docker

test: $(MARKERS)/test_go

analyze: analyze_go analyze_fossa

publish: gh-create-draft-release
	if [ "$(RELEASABLE)" = "yes" ]; then \
	  $(GH_RELEASE) upload $(VERSION) target/* ; \
	  for f in $$(find target -type f -exec basename {} \;); do \
	    $(TOOLCHAIN) aws s3 cp /project/target/$$f $(S3_BASE_URL)/$(VERSION)/$$f; \
	  done ; \
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
