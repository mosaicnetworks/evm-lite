NODES = 4
VALIDATORS = 4
POA = true
DEST =  $(CURDIR)/conf
PASS = $(CURDIR)/../eth/pwd.txt
IPBASE = node
IPADD = 0

conf:
	rm -rf conf
	./../eth/scripts/build-eth-conf.sh $(NODES) $(DEST) $(PASS) $(VALIDATORS) $(POA)
	./scripts/build-babble-conf.sh $(NODES) $(IPBASE) $(IPADD) $(DEST)
	./scripts/edit-toml.sh $(NODES) $(IPBASE) $(IPADD)


compile:
	./scripts/compile_pregenesis.sh 	

	
.PHONY: conf compile
