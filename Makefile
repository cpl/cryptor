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
	@echo "mode: atomic" > build/report.out

	@for dir in $$(ls); \
	do \
	if ls $$dir/*.go &> /dev/null; then \
		cd $$dir; \
		go test -coverprofile=coverage.out -race -v -parallel 4; \
		cat coverage.out | tail -n +2 >> ../build/report.out; \
		rm coverage.out; \
		cd ..; \
	fi \
	done;

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

install:
	go build -i -o $$GOBIN/cryptor -v -x cmd/cryptor-cli/*.go

test-cli:
	@make build && \
	echo "TEST CLI"

test:
	@make cover && make test-cli

docker:
	@docker build . -t cryptor

container:
	@docker run -p $(PORT):2000/udp -td cryptor; \
