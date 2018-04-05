OUTPUT = ansible-sandbox

CWD = $(shell pwd)
BASEDIR = $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
SRCDIR = $(BASEDIR)/src
STATICDIR = $(BASEDIR)/static

.PHONY: run clean build

%.go:

$(OUTPUT): $(SRCDIR)/*.go
	cd $(SRCDIR); \
	GOPATH=$(BASEDIR):$(SRCDIR)/vendor \
	go build -o $(BASEDIR)/ansible-sandbox

build: $(OUTPUT)

run: build
	$(BASEDIR)/$(OUTPUT) --static $(STATICDIR)

clean:
	-rm -f $(OUTPUT)
