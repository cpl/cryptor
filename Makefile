.PHONY: test clean


clean:
	@rm -rf build/



test:
	@mkdir -p build/
	@go test -coverprofile=build/cover.out -v -count=1 -parallel 8 ./...

build/cover.out: test

cover: build/cover.out
	@go tool cover -html=build/cover.out
