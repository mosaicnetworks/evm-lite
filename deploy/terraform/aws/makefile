nodes = 1
consensus = solo
conf = $(CURDIR)/../../conf/$(consensus)/conf

apply:
	rm -f ips.dat && \
	terraform apply -auto-approve \
	    -var-file=secret.tfvars \
		-var command=$(consensus) \
		-var nodes=$(nodes) \
		-var conf=$(conf)
destroy:
	terraform destroy -auto-approve -var-file=secret.tfvars

.PHONY: apply destroy