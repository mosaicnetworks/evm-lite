NODE = node4

deps:
	cd crowd-funding && npm install

deploy:
	$(MAKE) -C ../deploy prebuilt PREBUILT=babblepoademo && \
	$(MAKE) -C ../deploy start CONSENSUS=babble NODES=10

stop:
	$(MAKE) -C ../deploy stop

watch:
	./scripts/watch.sh

crowd-funding-demo:
	./scripts/crowd-funding-demo.sh

add-node:
	./scripts/write_evmlc_config.sh node0 && \
	./scripts/nominate_node.sh node0 $(NODE) && \
	./scripts/all_vote_node.sh $(NODE) && \
	./scripts/set_docker_peers.sh $(NODE) && \
	./scripts/restart_docker_node.sh $(NODE)

.PHONY: deps deploy stop watch crowd-funding-demo add-node