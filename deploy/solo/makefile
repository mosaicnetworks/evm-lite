nodes = 4
dest =  $(CURDIR)/conf
pass = $(CURDIR)/../scripts/pwd.txt

conf: 
	rm -rf conf
	../scripts/build-eth-conf.sh $(nodes) $(dest) $(pass)

.PHONY: conf
	
