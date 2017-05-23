.PHONY: default clean clobber update build install preinstall test run
.PHONY: killer lint test unit-test cov-view setup uncovered
.PHONY: fix-readme

default: install

export MYAPP := apidCRUD
export COV_DIR := cov
export COV_FILE := $(COV_DIR)/covdata.out
export COV_HTML := $(COV_DIR)/$(MYAPP)-coverage.html
export LOG_DIR := logs
export UNIT_TEST_DB := unit-test.db
VENDOR_DIR := github.com/30x/$(MYAPP)/vendor
SQLITE_PKG := github.com/mattn/go-sqlite3

clean:
	/bin/rm -f gen_*.go
	go clean
	/bin/rm -rf $(LOG_DIR)
	mkdir -p $(LOG_DIR)
	/bin/rm -rf $(COV_DIR)
	mkdir -p $(COV_DIR)
	/bin/rm -f $(UNIT_TEST_DB)

clobber: clean
	/bin/rm -rf ./vendor

update:
	glide --debug update

get:
	[ -d ./vendor ] \
	|| glide install

build: gen_swag.go
	time go $@

setup:
	mkdir -p $(LOG_DIR) $(COV_DIR)

# install this separately to speed up compilations.  thanks to Scott Ganyo.
preinstall: get
	[ -d $(VENDOR_DIR)/$(SQLITE_PKG) ] \
	|| go install $(VENDOR_DIR)/$(SQLITE_PKG)

install: setup preinstall gen_swag.go
	go $@ ./cmd/$(MYAPP)

run: install
	./runner.sh

killer:
	-pkill -f $(MYAPP)

test: unit-test

unit-test: gen_swag.go
	./unit-test.sh

cov-view:
	go tool cover -html=$(COV_FILE) -o $(COV_HTML)

func-test:
	./func-test.sh

lint: setup
	gometalinter.v1 --sort=path \
		-e "don't use underscores" \
		-e "should be" \
	| tee $(LOG_DIR)/$@.out

# not yet implemented
doc:
	godoc

fix-readme:
	./fix-readme.sh README.md template.txt

gen_swag.go: swagger.yaml cmd/swag/main.go
	go run cmd/swag/main.go $< > $@

uncovered:
	./uncovered.sh > $(LOG_DIR)/$@.out

##
## this is not a working target.
## these commands may not work on all platforms.
## these commands are only representative
## of what's needed, beyond go, to build this project.
##
#tools:
#	sudo apt-get install sqlite3
#	sudo apt-get install jq
#	go get -u gopkg.in/alecthomas/gometalinter.v1
#	gometalinter.v1 --install
#	go get github.com/Masterminds/glide
#	go get github.com/mattn/goveralls
