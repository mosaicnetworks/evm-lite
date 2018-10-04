EVM_LITE_VERSION=0.1.0
GIT_COMMIT=$(shell git rev-parse --short HEAD)

all: docker ami 
	
docker:
	go build \
		--ldflags '-extldflags "-static"' \
		-o ./docker/evml ./../../cmd/evml/ 
	docker build -t mosaicnetworks/evm-lite:$(EVM_LITE_VERSION)-$(GIT_COMMIT) \
		-t mosaicnetworks/evm-lite:latest \
		./docker/

ami:
	go build -o ami/evml ./../../cmd/evml/ && \
	packer build --var-file=ami/secret.json \
				 --var 'evml_version=$(EVM_LITE_VERSION)_$(GIT_COMMIT)' \
				  ami/template.json 

.PHONY: all docker ami