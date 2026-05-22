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

kube-apply:
	kubectl apply -f scripts/k8s/

kube-destroy:
	kubectl delete -f scripts/k8s/ --ignore-not-found

kube-restart: kube-destroy kube-apply

terraform:
	until nc -z 192.168.49.2 30002; do echo waiting for localstack; sleep 2; done;
	terraform -chdir=scripts/terraform/ plan
	terraform -chdir=scripts/terraform/ apply -auto-approve

terraform-init:
	terraform -chdir=scripts/terraform/ init

terraform-destroy:
	terraform -chdir=scripts/terraform/ destroy

minikube: kube-secrets kube-create terraform

test:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out