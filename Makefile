api = ledger-config
repository = clodoaldomarques

up: 
	docker compose up -d 
	$(MAKE) terraform

down: 
	docker compose down -v

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

restart: down up	

logs:
	docker compose logs

terraform:
	@if [ ! -d "scripts/terraform/.terraform" ]; then \
		echo "▶️  Inicializando Terraform..."; \
		terraform -chdir=scripts/terraform/ init; \
	else \
		echo "✅ Terraform já inicializado (pulando init)."; \
	fi
	@echo "⏳ Aguardando LocalStack na porta 4566..."
	@until nc -z localhost 4566; do echo "⏳ esperando..."; sleep 2; done
	@echo "📋 Gerando plano..."
	terraform -chdir=scripts/terraform/ plan
	@echo "🚀 Aplicando..."
	terraform -chdir=scripts/terraform/ apply -auto-approve

test:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out