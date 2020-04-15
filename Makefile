TARGETS = proto


CMD_DIR := ./cmd
PKG_DIR := ./pkg
OUT_DIR := ./out

DIST_TAR := dist.tar.gz
COV_FILE := cover.out


DIST_OS   += linux darwin
DIST_ARCH += amd64

GO_CMD ?= go
CGO_ENABLED ?= 0

EXTERNAL_TOOLS=\
	github.com/client9/misspell/cmd/misspell \
	github.com/golangci/golangci-lint/cmd/golangci-lint

GO111MODULE := on


GO_TEST_FLAGS := -v -count=1 -race -coverprofile=$(OUT_DIR)/$(COV_FILE) -covermode=atomic -parallel=16

.PHONY: deps
.PHONY: $(OUT_DIR) clean mod purge
.PHONY: test cover test-deps bench
.PHONY: lint-fmt lint-vet lint-misspell lint-ci lint
.PHONY: build dist


all: clean mod lint-fmt lint-vet test build

deps:
	@for tool in  $(EXTERNAL_TOOLS) ; do \
		echo "Installing/Updating $$tool" ; \
		GO111MODULE=off $(GO_CMD) get -u $$tool; \
	done

$(OUT_DIR):
	@mkdir -p $(OUT_DIR)

clean:
	@rm -rf $(OUT_DIR)

purge: clean
	$(GO_CMD) mod tidy
	$(GO_CMD) clean -cache
	$(GO_CMD) clean -testcache
	$(GO_CMD) clean -modcache

dist: clean $(OUT_DIR)
	@$(foreach arch, $(DIST_ARCH),			\
		$(foreach os,$(DIST_OS),			\
			$(foreach target, $(TARGETS),	\
			GOOS=$(os) GOARCH=$(arch) CGO_ENABLED=$(CGO_ENABLED) $(GO_CMD) build -o $(OUT_DIR)/$(target)_$(os)_$(arch) $(CMD_DIR)/$(target)/*.go && \
			shasum $(OUT_DIR)/$(target)_$(os)_$(arch) > $(OUT_DIR)/$(target)_$(os)_$(arch).shasum; \
		)))
	@cd $(OUT_DIR) && \
	tar -czvf $(DIST_TAR) --exclude $(DIST_TAR) *

build: $(OUT_DIR)
	$(foreach target,$(TARGETS),go build -o $(OUT_DIR)/$(target) $(CMD_DIR)/$(target)/*.go;)

test: $(OUT_DIR)
	$(GO_CMD) test $(GO_TEST_FLAGS) ./...

mod:
	$(GO_CMD) mod tidy
	$(GO_CMD) mod verify

cover: $(OUT_DIR)
	$(GO_CMD) tool cover -html=$(OUT_DIR)/$(COV_FILE)

test-deps:
	$(GO_CMD) test all

lint: lint-fmt lint-vet lint-misspell lint-ci

lint-fmt:
	$(GO_CMD) fmt ./...

lint-vet:
	$(GO_CMD) vet ./...

lint-ci:
	golangci-lint run -v ./...

lint-misspell:
	misspell -error ./...

bench:
	$(GO_CMD) test -bench=. -benchmem -benchtime=10s ./...
