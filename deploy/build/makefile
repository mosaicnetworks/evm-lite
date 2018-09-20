EVM_LITE_VERSION=0.1.0

all: docker ami 
	
docker:
	go build \
		--ldflags '-extldflags "-static"' \
		-o docker/evml ../../cmd/evml/ 
	docker build -t mosaicnetworks/evm-lite:$(EVM_LITE_VERSION) \
		-t mosaicnetworks/evm-lite:latest \
		docker/

ami:
	go build -o ami/evml ../../cmd/evml/ && \
	packer build --var-file=ami/secret.json ami/template.json 

.PHONY: all docker ami