nodes = 4

up: conf start

images :
	cd docker && make

conf: 
	rm -rf conf
	./scripts/build-eth-conf.sh $(nodes)

start:
	./scripts/run-testnet.sh $(nodes)

stop: 
	./scripts/stop-testnet.sh

.PHONY: up images conf start stop
	
