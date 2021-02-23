SUBDIRS := $(wildcard */.)
all: test

test: $(SUBDIRS)
$(SUBDIRS):
		cd $@ && go test

.PHONY: test $(SUBDIRS)
