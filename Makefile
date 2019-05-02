.PHONY: test clean test-color cover fmt

export GO111MODULE := on

TEST_FLAGS := -coverprofile=build/cover.out -covermode=atomic -v -timeout 30s -count=1 -parallel 8

clean:
	rm -rf build/;
	rm -f $(shell find . -name cover.out -type f);

test:
	@mkdir -p build/
	@go test $(TEST_FLAGS) ./...

test-color:
	@mkdir -p build/
	@go test $(TEST_FLAGS) ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/'';


build/cover.out: test

cover: build/cover.out
	@go tool cover -html=build/cover.out

fmt:
	$(eval GOFMT_OUT := $(shell gofmt -l . 2>&1))
	@if [ "$(GOFMT_OUT)" ]; then \
		echo "gofmt err in:\n$(GOFMT_OUT)"; \
		exit 1; \
	fi