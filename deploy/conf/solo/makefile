nodes = 4
dest =  $(CURDIR)/conf
pass = $(CURDIR)/../eth/pwd.txt

conf: 
	rm -rf conf
	./../eth/scripts/build-eth-conf.sh $(nodes) $(dest) $(pass)

.PHONY: conf
	
