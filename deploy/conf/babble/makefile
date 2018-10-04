NODES = 4
DEST =  $(CURDIR)/conf
PASS = $(CURDIR)/../eth/pwd.txt
IPBASE = node
IPADD = 0

conf:
	rm -rf conf
	./../eth/scripts/build-eth-conf.sh $(NODES) $(DEST) $(PASS)
	./scripts/build-babble-conf.sh $(NODES) $(IPBASE) $(IPADD) $(DEST)
	./scripts/edit-toml.sh $(NODES) $(IPBASE) $(IPADD)
.PHONY: conf