.PHONY: test clean


clean:
	@rm -rf build/
	@rm -rf **/cover.out

test:
	@mkdir -p build/
	@go test -coverprofile=build/cover.out -v -timeout 10s -count=1 -parallel 8 ./...

build/cover.out: test

cover: build/cover.out
	@go tool cover -html=build/cover.out
