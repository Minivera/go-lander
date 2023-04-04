vendor: go.sum
	go mod vendor

mod-update:
	go get -u -m
	go mod tidy

mod-tidy:
	go mod tidy

.PHONY: mod-update mod-tidy

GOROOT=$(shell go env GOROOT)

update_glue:
	cp $(GOROOT)/misc/wasm/wasm_exec.js ./example/index.js

test: update_glue
	GOOS=js GOARCH=wasm go test -exec $(GOROOT)/misc/wasm/go_js_wasm_exec -v $$(go list ./... | grep -v example)

build-example: update_glue
	GOOS=js GOARCH=wasm go build -o example/main.wasm example/main.go example/fondant.go

serve-example: build-example
	go run example/serve.go

.PHONY: update_glue test build-test serve-test
