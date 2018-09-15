EVM_LITE_VERSION=0.1.0

all: evm-lite-image 

evm-lite-image:
	go build \
		--ldflags '-extldflags "-static"' \
		-o evml ../../cmd/evml/
	docker build -t mosaicnetworks/evm-lite:$(EVM_LITE_VERSION) \
		-t mosaicnetworks/evm-lite:latest \
		.
	
.PHONY: all evm-lite-image