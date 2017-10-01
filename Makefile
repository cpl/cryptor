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
