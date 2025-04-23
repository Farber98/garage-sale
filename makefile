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

# ==============================================================================
# Modules support

run:
	go run api/services/sales/main.go | go run api/tooling/logfmt/main.go

tidy: 
	go mod tidy
	go mod vendor