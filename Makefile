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

build-test: test/fondant.go test/main.go
	GOOS=js GOARCH=wasm go build -o test/main.wasm test/main.go test/fondant.go

serve-test:
	go run test/serve.go

.PHONY: build-test serve-test
