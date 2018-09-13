consensus = solo
nodes = 4

up: images conf start

images :
	$(MAKE) -C docker

conf :
	$(MAKE) -C $(consensus) conf "nodes=$(nodes)"

start : 
	$(MAKE) -C terraform apply consensus=$(consensus) nodes=$(nodes)

stop:
	$(MAKE) -C terraform destroy

.PHONY: up images conf start stop
	
