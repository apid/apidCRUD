.PHONY: default clean build install preinstall test run
.PHONY: killer lint test unit-test coverage setup

default: install

MYAPP := apidCRUD
VENDOR_DIR := github.com/30x/$(MYAPP)/vendor
COV_DIR := cov
LOG_DIR := logs
SQLITE_PKG := github.com/mattn/go-sqlite3

clean:
	go clean
	/bin/rm -rf ./vendor
	/bin/rm -rf $(LOG_DIR)
	mkdir -p $(LOG_DIR)
	/bin/rm -rf $(COV_DIR)
	mkdir -p $(COV_DIR)

get:
	[ -d ./vendor ] \
	|| glide install

build:
	time go $@

setup:
	mkdir -p $(LOG_DIR) $(COV_DIR)

# install this separately to speed up compilations.  thanks to Scott Ganyo.
preinstall: get
	[ -d $(VENDOR_DIR)/$(SQLITE_PKG) ] \
	|| go install $(VENDOR_DIR)/$(SQLITE_PKG)

install: setup preinstall
	go $@ ./cmd/$(MYAPP)

run: install
	./runner.sh

killer:
	pkill -f $(MYAPP)

test: unit-test

unit-test:
	time go test | tee $(LOG_DIR)/$@.out

func-test:
	./tester.sh | tee $(LOG_DIR)/$@.out

coverage:
	./cover.sh | tee $(LOG_DIR)/$@.out

lint:
	gometalinter.v1 --sort=path -e "don't use underscores" \
	| tee $(LOG_DIR)/$@.out
