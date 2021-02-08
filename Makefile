.PHONY: help
help: ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/:.*##/:##/' | column -t -s '##'

.PHONY: test
test: ## Runs all the test
	@go test ./...

.PHONY: benchmark
benchmark: ## Runs the benchmark
	@go test ./... -bench=. -benchmem
