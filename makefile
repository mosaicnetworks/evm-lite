BUILD_TAGS?=evml

# vendor uses Glide to install all the Go dependencies in vendor/
vendor:
	glide install

test:
	glide novendor | xargs go test -count=1 -tags=unit

flagtest:
	glide novendor | xargs go test -count=1 -run TestFlagEmpty

.PHONY: vendor test
