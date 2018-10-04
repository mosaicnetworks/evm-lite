ENV = local # local or aws
VERSION = latest # evm-lite Docker Image
USER = 1000 # user that run evml inside the docker containers. (try 501 in MacOS)
CONSENSUS = solo # solo, babble, or raft
NODES = 1
IPBASE = node # use 10.0.2. for aws 
IPADD = 0 # use 10 for aws

up: conf start

build-docker-image:
	$(MAKE) -C build docker

build-ami:
	$(MAKE) -C build ami

conf:
	$(MAKE) -C conf/$(CONSENSUS) conf NODES=$(NODES) IPBASE=$(IPBASE) IPADD=$(IPADD)

start: 
	$(MAKE) -C terraform/$(ENV) apply NODES=$(NODES) CONSENSUS=$(CONSENSUS) VERSION=$(VERSION) USER=$(USER)

stop:
	$(MAKE) -C terraform/$(ENV) destroy

.PHONY: up build-docker build-ami conf start stop
	
