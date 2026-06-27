api = ledger-config
repository = clodoaldomarques

run:
	go run cmd/main.go

build:
	docker build -t $(repository)/$(api):$(version) -f scripts/docker/api/Dockerfile .
	docker tag $(repository)/$(api):$(version) $(repository)/$(api):latest

push:
	docker push $(repository)/$(api):$(version)
	docker push $(repository)/$(api):latest

publish: build push

version:
	docker images | grep $(api)

apply: 
	kubectl apply -f scripts/k8s/
	until nc -z 192.168.49.2 30002; do echo waiting for localstack; sleep 2; done;
	terraform -chdir=scripts/terraform/ plan
	terraform -chdir=scripts/terraform/ apply -auto-approve

destroy:
	kubectl delete -f scripts/k8s/ --ignore-not-found
	terraform -chdir=scripts/terraform/ destroy -auto-approve

restart: destroy apply

terraform:
	terraform -chdir=scripts/terraform/ init

test:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out