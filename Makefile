.PHONY: cover view push update build

profile-cpu:
	@cd $(PKG); \
	go test -race -parallel 4 -cpuprofile prof.cpu; \
	go tool pprof $(PKG).test ./prof.cpu; \
	rm prof.cpu $(PKG).test; \

profile-mem:
	@cd $(PKG); \
	go test -race -parallel 4 -memprofile prof.mem; \
	go tool pprof $(PKG).test ./prof.mem; \
	rm prof.mem $(PKG).test; \

cover:
	@mkdir -p build
	@CRYPTORROOT=`pwd`;
	@for pkg in `go list ./...`; do \
		cd $$GOPATH/src;\
		cd $$pkg; \
		go test -coverprofile=coverage.out -v -race -parallel 8; \
	done; \
	cd $$CRYPTORROOT;
	@gocovmerge $(shell find . -name coverage.out -type f) > build/report.out
	@rm -f $(shell find . -name coverage.out -type f);

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

build:
	go build -o build/cryptor -v -x cmd/cryptor-cli/*.go

install-cli:
	go build -i -o $$GOBIN/cryptor-cli -v -x cmd/cryptor-cli/*.go

test-cli:
	@make build && \
	echo "TEST CLI"

test:
	@make cover && make test-cli

docker:
	@docker build . -t cryptor

container:
	@docker run -p $(PORT):2000/udp -td cryptor; \

bench:
	@go test -bench=. ./...