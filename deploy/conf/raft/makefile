nodes = 3
dest =  $(CURDIR)/conf
pass = $(CURDIR)/../eth/pwd.txt
ipbase = node
ipadd = 0

conf:
	rm -rf conf
	./../eth/scripts/build-eth-conf.sh $(nodes) $(dest) $(pass)
	./scripts/build-raft-conf.sh $(nodes) $(ipbase) $(ipadd)
	./scripts/edit-toml.sh $(nodes) $(ipbase) $(ipadd)

.PHONY: conf