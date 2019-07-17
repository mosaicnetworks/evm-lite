BUILD_TAGS?=evml

# vendor uses Glide to install all the Go dependencies in vendor/
vendor:
	glide install

test:
	glide novendor | xargs go test

.PHONY: vendor test
