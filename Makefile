#!/usr/bin/make -f

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
SDK_PACK := $(shell go list -m github.com/bitcv-chain/bitcv-chain | sed  's/ /\@/g')

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/bitcv-chain/bitcv-chain/version.Name=bac \
		  -X github.com/bitcv-chain/bitcv-chain/version.ServerName=bacd \
		  -X github.com/bitcv-chain/bitcv-chain/version.ClientName=baccli \
		  -X github.com/bitcv-chain/bitcv-chain/version.Version=$(VERSION) \
		  -X github.com/bitcv-chain/bitcv-chain/version.Commit=$(COMMIT)"
ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/bitcv-chain/bitcv-chain/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

# The below include contains the tools target.
include contrib/devtools/Makefile

all: install lint check

build: go.sum
ifeq ($(OS),Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/bacd.exe ./cmd/bacd
	go build -mod=readonly $(BUILD_FLAGS) -o build/baccli.exe ./cmd/baccli
else
	go build  $(BUILD_FLAGS) -o build/bacd ./cmd/bacd
	go build -mod=readonly $(BUILD_FLAGS) -o build/baccli ./cmd/baccli
endif

build_bacdebug: go.sum
ifeq ($(OS),Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/bacdebug.exe ./cmd/bacdebug
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/bacdebug ./cmd/bacdebug
endif


build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

build-contract-tests-hooks:
ifeq ($(OS),Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/contract_tests.exe ./cmd/contract_tests
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/contract_tests ./cmd/contract_tests
endif

install: go.sum check-ledger
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/bacd
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/baccli

install_bacd: go.sum check-ledger
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/bacd

install_baccli: go.sum check-ledger
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/baccli

install-debug: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/bacdebug



########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/bacd -d 2 | dot -Tpng -o dependency-graph.png

clean:
	rm -rf snapcraft-local.yaml build/

distclean: clean
	rm -rf vendor/

########################################
### Testing


check: check-unit check-build
check-all: check check-race check-cover

check-unit:
	@VERSION=$(VERSION) go test -v -mod=readonly -tags='ledger test_ledger_mock' ./...

check-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock' ./...

check-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

check-build: build
	@go test -mod=readonly -p 4 `go list ./cli_test/...` -tags=cli_test


lint: ci-lint
ci-lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs goimports -w -local github.com/bitcv-chain/bitcv-chain

benchmark:
	@go test -mod=readonly -bench=. ./...


########################################
### Local validator nodes using docker and docker-compose

build-docker-bacdnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/bacd/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/bacd:Z tendermint/bacdnode testnet --v 4 -o . --starting-ip-address 192.168.10.2 ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down

setup-contract-tests-data:
	echo 'Prepare data for the contract tests'
	rm -rf /tmp/contract_tests ; \
	mkdir /tmp/contract_tests ; \
	cp "${GOPATH}/pkg/mod/${SDK_PACK}/client/lcd/swagger-ui/swagger.yaml" /tmp/contract_tests/swagger.yaml ; \
	./build/bacd init --home /tmp/contract_tests/.bacd --chain-id lcd contract-tests ; \
	tar -xzf lcd_test/testdata/state.tar.gz -C /tmp/contract_tests/

start-bac: setup-contract-tests-data
	./build/bacd --home /tmp/contract_tests/.bacd start &
	@sleep 2s

setup-transactions: start-bac
	@bash ./lcd_test/testdata/setup.sh

run-lcd-contract-tests:
	@echo "Running Bac LCD for contract tests"
	./build/baccli rest-server --laddr tcp://0.0.0.0:8080 --home /tmp/contract_tests/.baccli --node http://localhost:26657 --chain-id lcd --trust-node true

contract-tests: setup-transactions
	@echo "Running Bac LCD for contract tests"
	dredd && pkill bacd

# include simulations
include sims.mk

.PHONY: all build-linux install install-debug \
	go-mod-cache draw-deps clean build \
	setup-transactions setup-contract-tests-data start-bac run-lcd-contract-tests contract-tests \
	check check-all check-build check-cover check-ledger check-unit check-race
