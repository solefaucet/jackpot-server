#! /usr/bin/make
#
# Makefile for jackpot-server
#
# Targets:
# - "depend" retrieves the Go packages needed to run the linter and tests
# - "metalint" runs the linter and checks the code format using goimports
# - "test" runs the tests
#
# Meta targets:
# - "all" is the default target, it runs all the targets in the order above.
#
DEPEND=\
			 golang.org/x/tools/cmd/cover \
			 bitbucket.org/liamstask/goose/cmd/goose \
			 github.com/alecthomas/gometalinter

all: depend metalint test

depend:
	@go get -v $(DEPEND)

metalint:
	gometalinter \
		--deadline=60s \
		--disable-all \
		--vendor \
		--enable=goimports \
		--enable=golint \
		--enable=vetshadow \
		--enable=goconst \
		--enable=gosimple \
		--enable=staticcheck \
		--enable=dupl \
		--enable=gocyclo \
		--linter='dupl:dupl -plumbing -threshold {duplthreshold} ./*.go | grep -v "_test.go"::(?P<path>[^\s][^:]+?\.go):(?P<line>\d+)-\d+:\s*(?P<message>.*)' \
		./...

test:
	# run test
	go test -cover ./...
	# cleanup
	mysql -uroot -e 'drop database if exists jackpot_test;'
