SUBDIRS := $(wildcard */.)

all: $(SUBDIRS)
$(SUBDIRS):
		cd $@ && go test

.PHONY: all $(SUBDIRS)
