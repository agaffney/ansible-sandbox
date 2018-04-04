.PHONY: run

OUTPUT = ansible-sandbox

CWD = $(shell pwd)
BASEDIR = $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
SRCDIR = $(BASEDIR)/src
STATICDIR = $(BASEDIR)/static

%.go:

$(OUTPUT): $(SRCDIR)/*.go
	cd $(SRCDIR); \
	GOPATH=$(BASEDIR):$(SRCDIR)/vendor \
	go build -o $(BASEDIR)/ansible-sandbox

run: $(OUTPUT)
	$(BASEDIR)/$(OUTPUT) --static $(STATICDIR)
