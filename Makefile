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

version:
	docker images | grep $(api)


apply: terraform
	kubectl apply -f scripts/k8s/


destroy:
	kubectl delete -f scripts/k8s/ --ignore-not-found
	terraform -chdir=scripts/terraform/ destroy -auto-approve

reload: destroy apply


# =============================================================================
# Testes
# =============================================================================

test:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out

# =============================================================================
# Benchmark
# =============================================================================

# roda todos os benchmarks, uma vez, com alocação
bench:
	go test -bench=. -benchmem -run=^$$ ./...

# roda com múltiplas amostras (bom para comparar com benchstat)
bench-count:
	go test -bench=. -benchmem -run=^$$ -count=10 ./...

# roda só um benchmark específico (uso: make bench-one name=BenchmarkCreateEvent_Success)
bench-one:
	go test -bench=$(name) -benchmem -run=^$$ ./...

# salva baseline para comparação futura
bench-baseline:
	go test -bench=. -benchmem -run=^$$ -count=10 ./... | tee baseline.txt

# roda de novo depois de otimizar e compara com benchstat
bench-compare:
	go test -bench=. -benchmem -run=^$$ -count=10 ./... | tee after.txt
	benchstat baseline.txt after.txt

# detecta contenção sob concorrência
bench-parallel:
	go test -bench=Parallel -benchmem -run=^$$ -cpu=1,2,4,8 ./...

# =============================================================================
# Profiling
# =============================================================================

# CPU profile do caminho feliz
profile-cpu:
	go test -bench=BenchmarkCreateEvent_Success$$ -run=^$$ -cpuprofile=cpu.prof -benchtime=5s ./internal/application/ledger
	go tool pprof -http=:8080 cpu.prof

# memória: total alocado (acha pontos de alocação)
profile-mem:
	go test -bench=BenchmarkCreateEvent_Success$$ -run=^$$ -memprofile=mem.prof -benchtime=5s ./internal/application/ledger
	go tool pprof -http=:8080 mem.prof

# memória em uso (acha vazamento)
profile-mem-inuse:
	go tool pprof -inuse_space mem.prof

# goroutine profile (útil se estiver rodando como serviço)
profile-goroutine:
	go tool pprof http://localhost:6060/debug/pprof/goroutine

# trace de execução (investiga escalonamento, bloqueios)
profile-trace:
	go test -bench=BenchmarkCreateEvent_Parallel$$ -run=^$$ -trace=trace.out -benchtime=5s ./internal/application/ledger
	go tool trace trace.out

# limpa artefatos de profiling
profile-clean:
	rm -f cpu.prof mem.prof trace.out baseline.txt after.txt