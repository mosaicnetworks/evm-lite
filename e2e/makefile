NODE = node4

deps:
	cd crowd-funding && npm install

deploy:
	$(MAKE) -C ../deploy check  && \
	$(MAKE) -C ../deploy prebuilt PREBUILT=babblepoademo && \
	$(MAKE) -C ../deploy startnoerr CONSENSUS=babble NODES=10 && \
        ./scripts/rebuild_ipsdat.sh && \
	./scripts/get_running_nodes.sh


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
	./scripts/restart_docker_node.sh $(NODE) && \
	./scripts/rebuild_ipsdat.sh

.PHONY: deps deploy stop watch crowd-funding-demo add-node
