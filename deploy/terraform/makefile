consensus = "solo"

apply:
	rm -f ips.dat
	terraform apply -auto-approve -var "command=$(consensus)" -var "nodes=$(nodes)"

destroy:
	terraform destroy -auto-approve

