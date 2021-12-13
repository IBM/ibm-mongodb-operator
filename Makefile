#!/bin/bash
#
# Copyright 2021 IBM Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.DEFAULT_GOAL:=help

OPERATOR_SDK ?= $(shell which operator-sdk)
CONTROLLER_GEN ?= $(shell which controller-gen)
KUSTOMIZE ?= $(shell which kustomize)
# Specify whether this repo is build locally or not, default values is '1';
# If set to 1, then you need to also set 'DOCKER_USERNAME' and 'DOCKER_PASSWORD'
# environment variables before build the repo.
BUILD_LOCALLY ?= 1
TARGET_GOOS=linux
TARGET_GOARCH=amd64

# The namespcethat operator will be deployed in
NAMESPACE=ibm-common-services

# Image URL to use all building/pushing image targets;
# Use your own docker registry and image name for dev/test by overridding the IMG and REGISTRY environment variable.
IMG ?= ibm-mongodb-operator
BUNDLE_IMAGE_NAME=ibm-mongodb-operator-bundle
REGISTRY ?= "hyc-cloud-private-integration-docker-local.artifactory.swg-devops.com/ibmcom"
CSV_VERSION ?= $(shell cat version/version.go | grep = | cut -d '"' -f2)

# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

QUAY_USERNAME ?=
QUAY_PASSWORD ?=

MARKDOWN_LINT_WHITELIST=https://quay.io/cnr
VCS_URL ?= https://github.com/IBM/ibm-mongodb-operator
VCS_REF ?= $(shell git rev-parse HEAD)
TESTARGS_DEFAULT := "-v"
export TESTARGS ?= $(TESTARGS_DEFAULT)
VERSION ?= $(shell cat ./version/version.go | grep "Version =" | awk '{ print $$3}' | tr -d '"')

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

LOCAL_OS := $(shell uname)
LOCAL_ARCH := $(shell uname -m)
ifeq ($(LOCAL_OS),Linux)
    TARGET_OS ?= linux
    XARGS_FLAGS="-r"
    STRIP_FLAGS=
else ifeq ($(LOCAL_OS),Darwin)
    TARGET_OS ?= darwin
    XARGS_FLAGS=
    STRIP_FLAGS="-x"
else
    $(error "This system's OS $(LOCAL_OS) isn't recognized/supported")
endif

include common/Makefile.common.mk

ifeq ($(BUILD_LOCALLY),0)
    export CONFIG_DOCKER_TARGET = config-docker
endif


##@ Development

check: lint-all ## Check all files lint error

code-dev: ## Run the default dev commands which are the go tidy, fmt, vet then execute the $ make code-gen
	@echo Running the common required commands for developments purposes
	- make code-tidy
	- make code-fmt
	- make code-vet
	@echo Running the common required commands for code delivery
	- make check
	- make test

manager: generate code-fmt code-vet ## Build manager binary
	go build -o bin/manager main.go

run: generate code-fmt code-vet manifests ## Run against the configured Kubernetes cluster in ~/.kube/config
	WATCH_NAMESPACE=ibm-common-services go run ./main.go

install: manifests kustomize ## Install CRDs into a cluster
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ## Uninstall CRDs from a cluster
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller in the configured Kubernetes cluster in ~/.kube/config
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMAGE_REPO)/$(OPERATOR_IMAGE_NAME):$(CSV_VERSION)
	$(KUSTOMIZE) build config/default | kubectl apply -f -

##@ Generate code and manifests

manifests: controller-gen ## Generate manifests e.g. CRD, RBAC etc.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=ibm-mongodb-operator webhook paths="./..." output:crd:artifacts:config=config/crd/bases

generate: controller-gen ## Generate code e.g. API etc.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

bundle-manifests:
	$(KUSTOMIZE) build config/manifests | $(OPERATOR_SDK) generate bundle \
	-q --overwrite --version $(CSV_VERSION) $(BUNDLE_METADATA_OPTS)
	$(OPERATOR_SDK) bundle validate ./bundle

generate-all: manifests kustomize operator-sdk ## Generate bundle manifests, metadata and package manifests
	$(OPERATOR_SDK) generate kustomize manifests -q
	- make bundle-manifests CHANNELS=beta,stable-v1 DEFAULT_CHANNEL=stable-v1

##@ Build

build:
	@echo "Building the ibm-mongodb-operator binary"
	@CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a -o manager main.go

build-test-image: code-dev
	docker build -t quay.io/$(QUAY_USERNAME)/my-mongodb-operator:test \
	--build-arg VCS_REF=$(VCS_REF) --build-arg VCS_URL=$(VCS_URL) \
	--build-arg GOARCH="amd64" -f Dockerfile .

build-image: $(CONFIG_DOCKER_TARGET)
	$(eval ARCH := $(shell uname -m|sed 's/x86_64/amd64/'))
	docker build -t $(REGISTRY)/$(IMG)-$(ARCH):$(VERSION) \
	--build-arg VCS_REF=$(VCS_REF) --build-arg VCS_URL=$(VCS_URL) \
	--build-arg GOARCH="amd64" -f Dockerfile .
	@if [ $(BUILD_LOCALLY) -ne 1 ] && [ "$(ARCH)" = "amd64" ]; then docker push $(REGISTRY)/$(IMG)-$(ARCH):$(VERSION); fi

build-bundle-image: ## Build the operator bundle image.
	$(eval ARCH := $(shell uname -m|sed 's/x86_64/amd64/'))
	docker build -f bundle.Dockerfile -t $(REGISTRY)/$(BUNDLE_IMAGE_NAME)-$(ARCH):$(VERSION) .

# runs on amd64 machine
build-image-ppc64le: $(CONFIG_DOCKER_TARGET)
ifeq ($(LOCAL_OS),Linux)
ifeq ($(LOCAL_ARCH),x86_64)
	docker run --rm --privileged multiarch/qemu-user-static:register --reset
	docker build -t $(REGISTRY)/$(IMG)-ppc64le:$(VERSION) \
	--build-arg VCS_REF=$(VCS_REF) --build-arg VCS_URL=$(VCS_URL) \
	--build-arg GOARCH="ppc64le" -f Dockerfile.ppc64le .
	@if [ $(BUILD_LOCALLY) -ne 1 ]; then docker push $(REGISTRY)/$(IMG)-ppc64le:$(VERSION); fi
endif
endif

# runs on amd64 machine
build-image-s390x: $(CONFIG_DOCKER_TARGET)
ifeq ($(LOCAL_OS),Linux)
ifeq ($(LOCAL_ARCH),x86_64)
	docker run --rm --privileged multiarch/qemu-user-static:register --reset
	docker build -t $(REGISTRY)/$(IMG)-s390x:$(VERSION) \
	--build-arg VCS_REF=$(VCS_REF) --build-arg VCS_URL=$(VCS_URL) \
	--build-arg GOARCH="s390x" -f Dockerfile.s390x .
	@if [ $(BUILD_LOCALLY) -ne 1 ]; then docker push $(REGISTRY)/$(IMG)-s390x:$(VERSION); fi
endif
endif

##@ Test

test: ## Run unit test
	@go test ${TESTARGS} ./controllers/...

test-e2e: ## Run integration e2e tests with different options.
	@echo ... Running the same e2e tests with different args ...
	@echo ... Running locally ...
	- operator-sdk test local ./test/e2e --verbose --up-local --namespace=${NAMESPACE}
	# @echo ... Running with the param ...
	# - operator-sdk test local ./test/e2e --namespace=${NAMESPACE}

##@ Release

images: build-image build-image-ppc64le build-image-s390x
ifeq ($(LOCAL_OS),Linux)
ifeq ($(LOCAL_ARCH),x86_64)
	@curl -L -o /tmp/manifest-tool https://github.com/estesp/manifest-tool/releases/download/v1.0.3/manifest-tool-linux-amd64
	@chmod +x /tmp/manifest-tool
	/tmp/manifest-tool push from-args --platforms linux/amd64,linux/ppc64le,linux/s390x --template $(REGISTRY)/$(IMG)-ARCH:$(VERSION) --target $(REGISTRY)/$(IMG) --ignore-missing
	/tmp/manifest-tool push from-args --platforms linux/amd64,linux/ppc64le,linux/s390x --template $(REGISTRY)/$(IMG)-ARCH:$(VERSION) --target $(REGISTRY)/$(IMG):$(VERSION) --ignore-missing
endif
endif

csv: ## Push CSV package to the catalog
	@RELEASE=${CSV_VERSION} common/scripts/push-csv.sh

remove-csv:
	common/scripts/delete-csv.sh

all: check test coverage build images

##@ Cleanup
clean: ## Clean build binary
	rm -f bin/manager

promote-to-beta:
	common/scripts/promote-to-beta.sh ${CSV_VERSION}

# download operator-sdk if necessary
operator-sdk:
ifeq (, $(OPERATOR_SDK))
	@./common/scripts/install-operator-sdk.sh
	OPERATOR_SDK=/usr/local/bin/operator-sdk
endif

# find or download kubebuilder
# download kubebuilder if necessary
kube-builder:
ifeq (, $(wildcard /usr/local/kubebuilder))
	@./common/scripts/install-kubebuilder.sh
endif

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

kustomize:
ifeq (, $(shell which kustomize))
	@{ \
	set -e ;\
	KUSTOMIZE_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$KUSTOMIZE_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/kustomize/kustomize/v3@v3.5.4 ;\
	rm -rf $$KUSTOMIZE_GEN_TMP_DIR ;\
	}
KUSTOMIZE=$(GOBIN)/kustomize
else
KUSTOMIZE=$(shell which kustomize)
endif

##@ Help
help: ## Display this help
	@echo "Usage:\n  make \033[36m<target>\033[0m"
	@awk 'BEGIN {FS = ":.*##"}; \
			/^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } \
			/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: all build run check install uninstall code-dev test test-e2e coverage images csv clean controller-gen kustomize help
