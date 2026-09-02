api = ledger-config
repository = clodoaldomarques

up: 
	docker compose up -d 
	$(MAKE) terraform

down: 
	docker compose down -v

run:
	export $$(cat .env | xargs) && go run cmd/main.go

build:
	docker build -t $(repository)/$(api):$(version) -f scripts/docker/api/Dockerfile .
	docker tag $(repository)/$(api):$(version) $(repository)/$(api):latest

push:
	docker push $(repository)/$(api):$(version)
	docker push $(repository)/$(api):latest

publish: build push

version:
	docker images | grep $(api)

restart: down up	

logs:
	docker compose logs

terraform:
	@if [ ! -d "scripts/terraform/.terraform" ]; then \
		terraform -chdir=scripts/terraform/ init;\
	fi
	until nc -z 192.168.67.2 30002; do echo waiting for localstack; sleep 2; done;
	terraform -chdir=scripts/terraform/ plan
	terraform -chdir=scripts/terraform/ apply -auto-approve

test:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out

version:
	docker images | grep $(api)


apply: 
	$(MAKE) terraform
	kubectl apply -f scripts/k8s/


destroy:
	kubectl delete -f scripts/k8s/ --ignore-not-found
	terraform -chdir=scripts/terraform/ destroy -auto-approve

reload: destroy apply