env = local # local or aws
version = 0.1.0 # evm-lite Docker Image
user = 1000 # user that run evml inside the docker containers. (try 501 in MacOS)
consensus = solo # solo, babble, or raft
nodes = 1
ipbase = node # use 10.0.2. for aws 
ipadd = 0 # use 10 for aws

up: conf start

build-docker-image:
	$(MAKE) -C build docker

build-ami:
	$(MAKE) -C build ami

conf:
	$(MAKE) -C conf/$(consensus) conf nodes=$(nodes) ipbase=$(ipbase) ipadd=$(ipadd)

start: 
	$(MAKE) -C terraform/$(env) apply nodes=$(nodes) consensus=$(consensus) version=$(version) user=$(user)

stop:
	$(MAKE) -C terraform/$(env) destroy

.PHONY: up build-docker build-ami conf start stop
	
