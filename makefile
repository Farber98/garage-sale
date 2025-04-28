# ==============================================================================
# Define deps and env vars

GOLANG          := golang:1.24
ALPINE          := alpine:3.21
KIND            := kindest/node:v1.32.0
POSTGRES        := postgres:17.4
GRAFANA         := grafana/grafana:11.6.0
PROMETHEUS      := prom/prometheus:v3.2.0
TEMPO           := grafana/tempo:2.7.0
LOKI            := grafana/loki:3.4.0
PROMTAIL        := grafana/promtail:3.4.0

KIND_CLUSTER    := garage-sale-cluster
NAMESPACE       := sales-system
SALES_APP       := sales
BASE_IMAGE_NAME := garage-sale
VERSION         := "0.0.1"
SALES_IMAGE     := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)


# ==============================================================================
# Install dependencies

dev-brew:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli
	brew list watch || brew install watch

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev-docker:
	docker pull $(GOLANG) & \
	docker pull $(ALPINE) & \
	docker pull $(KIND) & \
	docker pull $(POSTGRES) & \
	docker pull $(GRAFANA) & \
	docker pull $(PROMETHEUS) & \
	docker pull $(TEMPO) & \
	docker pull $(LOKI) & \
	docker pull $(PROMTAIL) & \
	wait;

# ==============================================================================
# Building containers

build: sales

sales:
	docker build \
		-f zarf/docker/dockerfile.sales \
		-t $(SALES_IMAGE) \
		--build-arg BUILD_REF=$(shell git rev-parse --short HEAD) \
		--build-arg BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# ==============================================================================
# Running from within k8s/kind

dev-up: 
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml
	
	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

	wait;

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status:
	watch -n 2 kubectl get pods -o wide --all-namespaces

dev-load:
	kind load docker-image $(SALES_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(SALES_APP) --timeout=120s --for=condition=Ready

# Restart POD. 
dev-restart:
	kubectl rollout restart deployment $(SALES_APP) --namespace=$(NAMESPACE)

dev-run: 
	build dev-up dev-load dev-apply

# When updating code we need to load and restart.
dev-update: build dev-load dev-restart

# When updating yaml, we need to load and apply.
dev-update-apply: build dev-load dev-apply

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(SALES_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run api/tooling/logfmt/main.go -service=$(SALES_APP)

# Info about Deployment of sales POD
dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(SALES_APP)

# Info about running Sales POD
dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(SALES_APP)

# ==============================================================================
# Metrics and Tracing

metrics-view:
	expvarmon -ports="localhost:3010" -endpoint="/debug/vars" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

statsviz:
	open http://localhost:3010/debug/statsviz

# ==============================================================================
# Modules support

run:
	go run api/services/sales/main.go | go run api/tooling/logfmt/main.go

tidy: 
	go mod tidy
	go mod vendor