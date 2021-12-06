SHELL := /bin/bash

# ==============================================================================
# Testing running system

# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

# For testing load on the service.
# hey -m GET -c 100 -n 10000 -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users/1/2

# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem

# ==============================================================================

run:
	go run app/services/wt-api/main.go --help | go run app/tooling/logfmt/main.go

admin:
	go run app/tooling/admin/main.go

tidy:
	go mod tidy
	go mod vendor

build:
	go build -ldflags "-X main.build=local"

# ==============================================================================
# Building containers

VERSION := 1.0

all: service

service:
	docker build \
		-f zarf/docker/Dockerfile \
		-t wt-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within k8s/kind

KIND_CLUSTER := wt-starter-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=wt-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	cd zarf/k8s/kind/wt-pod; kustomize edit set image wt-api-image=wt-api-amd64:$(VERSION)
	kind load docker-image wt-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/wt-pod | kubectl apply -f -

kind-update-apply: all kind-load kind-apply

kind-restart:
	kubectl rollout restart deployment wt-pod

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-logs:
	kubectl logs -l app=wt --all-containers=true -f --tail=100 --namespace=wt-system | go run app/tooling/logfmt/main.go

kind-status-service:
	kubectl get pods -o wide --watch --namespace=wt-system

kind-describe:
	kubectl describe nodes
	kubectl describe svc
	kubectl describe pod -l app=wt

kind-update: all kind-load kind-restart