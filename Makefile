.PHONY: binary run lint test help update_test_data

.DEFAULT_GOAL := help

all: test binary

binary: ## build binary for Linux
	./scripts/build/binary.sh

run: binary ## runs the newly created ./bin symlink
	./bin $(ARGS)

lint: ## run all the lint tools
	golint --set_exit_status .

test: ## run all tests
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

update_test_data:
	@go run update_data.go test-data

mongo_start: ## runs mongo
	docker run -d --name mongo -p 27017-27019:27017-27019 mongo --bind_ip_all

mongo_shell: ## access the mongo shell
	docker exec -it mongo mongo

docs: ## output router markdown docs to stdout
	./bin -doc

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
