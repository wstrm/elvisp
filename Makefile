SOURCEDIR=$(shell pwd)
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=elvisp

TESTPKGS=./ ./lease ./database

VERSION=`<./VERSION`
BUILD_TIME=`date +%FT%T%z`


LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

.PHONY: install deploy clean build init testshort testlong race

all: init build testshort
build: $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY}

init:
	go get github.com/Masterminds/glide
	glide install

install:
	go install ${LDFLAGS} ./...

deploy:
	go build ${LDFLAGS} -tags deploy -o ${BINARY}

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

testshort:
	go test -short ${TESTPKGS}

testlong:
	go test ${TESTPKGS}

race:
	go test -race ${TESTPKGS}
