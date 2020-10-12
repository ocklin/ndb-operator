# Copyright 2019 Oracle and/or its affiliates. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

VERSION ?= 1.0.0

ARCH     ?= amd64
OS       ?= darwin
UNAME_S  := $(shell uname -s)
PKG      := github.com/ocklin/ndb-operator/
CMD_DIRECTORIES := $(sort $(dir $(wildcard ./cmd/*/)))
COMMANDS := $(CMD_DIRECTORIES:./cmd/%/=%)

BASEDIR=/home/bo/prg/mysql-bld/trunk
RTDIR=${BASEDIR}/runtime_output_directory

BINDIR   :=bin/
SBINDIR  :=sbin/
LIBDIR   :=lib64/mysql/


ifeq ($(UNAME_S),Darwin)
ifeq ($(OS),linux)
	  # Cross-compiling from OSX to linux, go install puts the binaries in $GOPATH/bin/$GOOS_$GOARCH
    BINARIES := $(addprefix $(GOPATH)/bin/$(OS)_$(ARCH)/,$(COMMANDS))
else
	  # Compiling on darwin for linux, go install puts the binaries in $GOPATH/bin
    BINARIES := $(addprefix $(GOPATH)/bin/,$(COMMANDS))
endif
else
ifeq ($(UNAME_S),Linux)
	# Compiling on linux for linux, go install puts the binaries in $GOPATH/bin
    BINARIES := $(addprefix $(GOPATH)/bin/,$(COMMANDS))
else
	$(error "Unsupported OS: $(UNAME_S)")
endif
endif

.PHONY: all
all: build

.PHONY: build
build: 
	@echo "Building: $(BINARIES)"
	@echo "arch:     $(ARCH)"
	@echo "os:       $(OS)"
	@echo "version:  $(VERSION)"
	@echo "pkg:      $(PKG)"
	@echo "bin:      $(BINARIES)"
	@touch pkg/version/version.go # Important. Work around for https://github.com/golang/go/issues/18369
	ARCH=$(ARCH) OS=$(OS) VERSION=$(VERSION) PKG=$(PKG) ./hack/build.sh
	mkdir -p ./bin/$(OS)_$(ARCH)/
	cp $(BINARIES) ./bin/$(OS)_$(ARCH)/

.PHONY: binaries
binaries:
	@echo $(BINARIES)

# install-minimal 
#   copies the needed ndb binaries 
#
# from a build (not install) directory 
# to this local folder 
# for going into a container
.PHONY: install-minimal
install-minimal:
	install -m 0750 -d bin/mysql/$(SBINDIR)
	install -m 0750 -d bin/mysql/$(BINDIR)
	install -m 0750 -d bin/mysql/$(LIBDIR)

	install -m 0755 $(RTDIR)/mysqld \
					$(RTDIR)/mysqladmin \
					$(RTDIR)/ndb_mgmd \
					$(RTDIR)/ndbmtd bin/mysql/$(SBINDIR)

	install -m 0755 $(RTDIR)/mysql \
					$(RTDIR)/ndb_mgm bin/mysql/$(BINDIR)

# just a convenience as I never remember which way around
.PHONY: docker-build
docker-build: build-docker

.PHONY: build-docker
build-docker: build-docker-cluster 
	# build-docker-agent

.PHONY: build-docker-cluster
build-docker-cluster: install-minimal
	@docker build \
	-t mysql-cluster:$(VERSION) \
	-f docker/Dockerfile .

.PHONY: build-docker-agent
build-docker-agent:
	@docker build \
	-t ndb-agent:$(VERSION) \
	-f docker/ndb-agent/Dockerfile .


.PHONY: version
version:
	@echo $(VERSION)

.PHONY: clean
clean:
	rm -rf .go bin

.PHONY: generate
generate:
	./hack/update-codegen.sh

run:
	bin/linux_amd64/ndb-operator --kubeconfig=/home/bo/.kube/config 