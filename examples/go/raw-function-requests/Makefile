init::
	terraform -chdir=./tf init

apply::
	terraform -chdir=./tf apply -auto-approve

destroy::
	terraform -chdir=./tf destroy -auto-approve

test::
	curl -XPOST \
  "https://$$(terraform -chdir=./tf output -raw api_gateway_endpoint)/demo" \
  -d '{"name": "test"}' \
  -H "Content-Type: application/json"
