vendor: go.sum
	GO111MODULE=on go mod vendor

mod-update:
	GO111MODULE=on go get -u -m
	GO111MODULE=on go mod tidy

mod-tidy:
	GO111MODULE=on go mod tidy

.PHONY: $(CMD_PKGS)
.PHONY: mod-update mod-tidy

#################################################
# Test and linting
#################################################

test: vendor generated
	@CGO_ENABLED=0 go test -v $$(go list ./... | grep -v example)

comma:= ,
empty:=
space:= $(empty) $(empty)

COVER_TEST_PKGS:=$(shell find . -type f -name '*_test.go' | grep -v vendor | rev | cut -d "/" -f 2- | rev | grep -v example | sort -u)
$(COVER_TEST_PKGS:=-cover): %-cover: all-cover.txt
	@CGO_ENABLED=0 go test -coverprofile=$@.out -covermode=atomic ./$*
	@if [ -f $@.out ]; then \
		grep -v "mode: atomic" < $@.out >> all-cover.txt; \
		rm $@.out; \
	fi

all-cover.txt:
	echo "mode: atomic" > all-cover.txt

cover: vendor generated all-cover.txt $(COVER_TEST_PKGS:=-cover)

.golangci.gen.yml: .golangci.yml
	$(shell awk '/enable:/{y=1;next} y == 0 {print}' $< > $@)
LINTERS=$(filter-out megacheck,$(shell awk '/enable:/{y=1;next} y != 0 {print $$2}' .golangci.yml))

lint: vendor/bin/golangci-lint vendor .golangci.gen.yml
	$< run -c .golangci.gen.yml $(LINTERS:%=-E %) ./...
	$< run -c .golangci.gen.yml -E megacheck ./...
	$< run -c .golangci.gen.yml -E goimports ./...

.PHONY: lint

#################################################
# Building and serving examples
#################################################

build-example: test/fondant.go test/main.go
	GOOS=js GOARCH=wasm go build -o example/main.wasm example/main.go example/fondant.go

serve-example: build-example
	go run example/serve.go

.PHONY: build-test serve-test
