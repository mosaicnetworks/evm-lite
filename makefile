BUILD_TAGS?=evml

# vendor uses Glide to install all the Go dependencies in vendor/
vendor:
	(rm glide.lock || rm -rf vendor ) && glide install

install:
	go install ./cmd/evml/

test:
	glide novendor | xargs go test -count=1 -tags=unit

flagtest:
	glide novendor | xargs go test -count=1 -run TestFlagEmpty

 
lint:
	glide novendor | xargs golint

.PHONY: vendor test flagtest lint
