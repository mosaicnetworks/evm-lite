NODES = 4
DEST =  $(CURDIR)/conf
PASS = $(CURDIR)/../eth/pwd.txt

conf: 
	rm -rf conf
	./../eth/scripts/build-eth-conf.sh $(NODES) $(DEST) $(PASS)

compile: ;

.PHONY: conf compile
	
