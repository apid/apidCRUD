.PHONY: default clean build install preinstall test run
.PHONY: killer unit-test coverage setup

default: install

MYAPP := apidCRUD
VENDOR_DIR := github.com/30x/$(MYAPP)/vendor
COVDIR := cov
LOGDIR := logs
SQLITE_PKG := github.com/mattn/go-sqlite3

clean:
	go clean
	/bin/rm -rf ./vendor
	/bin/rm -rf $(LOGDIR)
	mkdir -p $(LOGDIR)
	/bin/rm -rf $(COVDIR)
	mkdir -p $(COVDIR)

get:
	[ -d ./vendor ] \
	|| glide install

build test:
	time go $@

setup:
	mkdir -p $(LOGDIR) $(COVDIR)

# install this separately to speed up compilations.  thanks to Scott Ganyo.
preinstall: get
	[ -d $(VENDOR_DIR)/$(SQLITE_PKG) ] \
	|| go install $(VENDOR_DIR)/$(SQLITE_PKG)

install: setup preinstall
	go $@ ./cmd/$(MYAPP)

run: install
	./runner.sh

killer:
	pkill -f apidCRUD

unit-test:
	time go test

func-test:
	./tester.sh

coverage:
	./cover.sh

lint:
	golint
