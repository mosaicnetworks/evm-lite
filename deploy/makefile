ENV = local # local or aws
VERSION = 0.2.0 # evm-lite Docker Image
CONSENSUS = solo # solo, babble, or raft
NODES = 1
IPBASE = node # use 10.0.2. for aws 
IPADD = 0 # use 10 for aws
PREBUILT = "" # prebuilt configurations

up: conf start

build-docker-image:
	$(MAKE) -C build docker

build-ami:
	$(MAKE) -C build ami

conf:
	$(MAKE) -C conf/$(CONSENSUS) conf NODES=$(NODES) IPBASE=$(IPBASE) IPADD=$(IPADD)

prebuilt:

	$(MAKE) -C prebuilt prebuilt PREBUILT="$(PREBUILT)"

start: 
	$(MAKE) -C terraform/$(ENV) apply NODES=$(NODES) CONSENSUS=$(CONSENSUS) VERSION=$(VERSION)

stop:
	$(MAKE) -C terraform/$(ENV) destroy

.PHONY: up build-docker build-ami conf prebuilt start stop
	
