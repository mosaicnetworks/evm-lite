NODES = 1
CONSENSUS = solo
CONF = $(CURDIR)/../../conf/$(CONSENSUS)/conf

apply:
	rm -f ips.dat && \
	terraform apply -auto-approve \
	    -var-file=secret.tfvars \
		-var command="run $(CONSENSUS)" \
		-var nodes=$(NODES) \
		-var conf=$(CONF)
destroy:
	terraform destroy -auto-approve -var-file=secret.tfvars

.PHONY: apply destroy