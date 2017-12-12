.PHONY: cover view push update build

profile-cpu:
	@cd $(PKG); \
	go test -race -parallel 4 -cpuprofile prof.cpu; \
	go tool pprof $(PKG).test ./prof.cpu; \

profile-mem:
	@cd $(PKG); \
	go test -race -parallel 4 -memprofile prof.mem; \
	go tool pprof $(PKG).test ./prof.mem; \

cover:
	gocovmerge $(shell find . -name coverage.out -type f) > build/report.out

view:
	@go tool cover -html=build/report.out

push:
	@if [ -n "$$(git status --porcelain)" ]; then \
		git status; \
	else \
		git push; \
	fi \

update:
	git pull

test:
	@mkdir -p build
	@CRYPTORROOT=`pwd`;
	@for pkg in `go list ./...`; do \
		cd $$GOPATH/src;\
		cd $$pkg; \
		go test -coverprofile=coverage.out -v -race -parallel 8 | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''; \
	done; \
	cd $$CRYPTORROOT;

docker:
	@docker build . -t cryptor

container:
	@docker run -p $(PORT):2000/udp -td cryptor; \

bench:
	@go test -bench=. ./...

clean:
	@rm -f $(shell find . -name coverage.out -type f);
	@rm -f prof.cpu $(PKG).test;
	@rm -f prof.mem $(PKG).test;
