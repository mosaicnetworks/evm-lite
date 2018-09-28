user = 1000
nodes = 1
consensus = solo
version = 0.1.0
conf = $(CURDIR)/../../conf/$(consensus)/conf

apply:
	rm -f ips.dat && \
	terraform apply -auto-approve \
		-var user=$(user) \
		-var command=$(consensus) \
		-var nodes=$(nodes) \
		-var version=$(version) \
		-var conf=$(conf)
destroy:
	terraform destroy -auto-approve

.PHONY: apply destroy