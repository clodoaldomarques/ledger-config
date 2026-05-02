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

kube-secrets:
	kubectl create secret generic aws-secrets --from-literal=AWS_ACCESS_KEY_ID='test' --from-literal=AWS_SECRET_ACCESS_KEY='test' --from-literal=AWS_SESSION_TOKEN='' --from-literal=aws-account='000000000000' --from-literal=aws-assume-role='' --from-literal=aws-region='us-east-1'

kube-create:
	kubectl apply -f scripts/k8s/localstack-service.yaml
	kubectl apply -f scripts/k8s/dynamodb-admin-service.yaml
	kubectl apply -f scripts/k8s/ollama-service.yaml
	kubectl apply -f scripts/k8s/app-service.yaml

kube-delete:
	kubectl delete -f scripts/k8s/localstack-service.yaml
	kubectl delete -f scripts/k8s/dynamodb-admin-service.yaml
	kubectl delete -f scripts/k8s/ollama-service.yaml
	kubectl delete -f scripts/k8s/app-service.yaml

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