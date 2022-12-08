# Update dependencies
.PHONY: update
update:
	go get -u ./...
	go mod tidy
	go mod tidy
