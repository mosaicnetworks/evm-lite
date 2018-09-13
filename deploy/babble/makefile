nodes = 4

conf:
	rm -rf conf
	./scripts/build-eth-conf.sh $(nodes)
	./scripts/build-babble-conf.sh $(nodes)
	./scripts/build-base-conf.sh $(nodes)

.PHONY: conf