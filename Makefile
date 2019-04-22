.PHONY: test clean


clean:
	@rm -rf build/
	@rm -rf **/cover.out

test:
	@mkdir -p build/
	@go test -coverprofile=build/cover.out -v -timeout 30s -count=1 -parallel 8 ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/'';

build/cover.out: test

cover: build/cover.out
	@go tool cover -html=build/cover.out
