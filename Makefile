.PHONY: default clean build install preinstall test run killer unit-test coverage

default: install

MYAPP := apidCRUD
VENDOR_DIR := github.com/30x/$(MYAPP)/vendor
COVDIR := cov
LOGDIR := logs

clean:
	go "$@"
	/bin/rm -rf $(LOGDIR)
	mkdir -p $(LOGDIR)
	/bin/rm -rf $(COVDIR)
	mkdir -p $(COVDIR)

build test:
	time go $@

# install this separately to speed up compilations.  thanks to Scott Ganyo.
preinstall:
	go install $(VENDOR_DIR)/github.com/mattn/go-sqlite3

install: preinstall
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
