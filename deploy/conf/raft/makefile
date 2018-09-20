nodes = 4
dest =  $(CURDIR)/conf
pass = $(CURDIR)/../scripts/pwd.txt

conf:
	rm -rf conf
	../scripts/build-eth-conf.sh $(nodes) $(dest) $(pass)
	./scripts/build-raft-conf.sh $(nodes)
	./scripts/edit-toml.sh $(nodes)

.PHONY: conf