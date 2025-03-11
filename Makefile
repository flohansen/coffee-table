.PHONY: setup
setup:
	cp ./.githooks/pre-push ./.git/hooks/

.PHONY: gen
gen:
	mkdir -p pkg/proto/
	go generate ./...

.PHONY: test
test: gen
	go test ./... -race -cover

.PHONY: clean
clean:
	rm -rf pkg/proto/ ./.git/hooks/pre-push
