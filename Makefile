.PHONY: default clean clobber update build install preinstall test run
.PHONY: killer lint test unit-test cov-view setup

default: install

export MYAPP := apidCRUD
export COV_DIR := cov
export COV_FILE := $(COV_DIR)/covdata.out
export COV_HTML := $(COV_DIR)/$(MYAPP)-coverage.html
export LOG_DIR := logs
VENDOR_DIR := github.com/30x/$(MYAPP)/vendor
SQLITE_PKG := github.com/mattn/go-sqlite3

clean:
	go clean
	/bin/rm -rf $(LOG_DIR)
	mkdir -p $(LOG_DIR)
	/bin/rm -rf $(COV_DIR)
	mkdir -p $(COV_DIR)

clobber: clean
	/bin/rm -rf ./vendor

update:
	glide --debug update

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
	-pkill -f $(MYAPP)

test: unit-test

unit-test:
	./unit-test.sh

cov-view:
	go tool cover -html=$(COV_FILE) -o $(COV_HTML)

func-test:
	./func-test.sh

lint:
	gometalinter.v1 --sort=path -e "don't use underscores" \
	| tee $(LOG_DIR)/$@.out

# not yet implemented
doc:
	godoc
