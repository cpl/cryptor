.PHONY: cover push update build test testf clean tool

profile-cpu:
	@cd $(PKG); \
	go test -race -parallel 4 -cpuprofile prof.cpu; \
	go tool pprof $(PKG).test ./prof.cpu; \

profile-mem:
	@cd $(PKG); \
	go test -race -parallel 4 -memprofile prof.mem; \
	go tool pprof $(PKG).test ./prof.mem; \

tool:
	@go build -o build/$(TARGET) tools/$(TARGET)/$(TARGET).go;
	@chmod a+x build/$(TARGET);
	@echo "DONE $(TARGET)";

cover:
	@go tool cover -html=build/report.out

update:
	git fetch --all
	git pull

push:
	@if [ -n "$$(git status --porcelain)" ]; then \
		git status; \
	else \
		git push; \
	fi \

test: clean
	@mkdir -p build
	@CRYPTORROOT=`pwd`;
	@cd $$CRYPTORROOT;
	@go test -coverprofile=build/report.out -v -count=1 -race -parallel 8 ./...; \

testf: clean
	@mkdir -p build;
	@CRYPTORROOT=`pwd`;
	@cd $$CRYPTORROOT;
	@go test -coverprofile=build/report.out -v -count=1 -race -parallel 8 ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''; \

testall: update clean testf bench

bench:
	@go test -bench=. ./...;
	cd $$CRYPTORROOT;

clean:
	@rm -f $(shell find . -name coverage.out -type f);
	@rm -f prof.cpu $(PKG).test;
	@rm -f prof.mem $(PKG).test;
