PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation' | grep -v '/prometheus' | grep -v '/clitest' | grep -v '/lcd' | grep -v '/protobuf')
PACKAGES_MODULES=$(shell go list ./... | grep 'modules')
PACKAGES_TYPES=$(shell go list ./... | grep 'irisnet/irishub/types')
PACKAGES_STORE=$(shell go list ./... | grep 'irisnet/irishub/store')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')

all: get_tools get_vendor_deps install

COMMIT_HASH := $(shell git rev-parse --short HEAD)

InvariantLevel := $(shell if [ -z ${InvariantLevel} ]; then echo "panic"; else echo ${InvariantLevel}; fi)
NetworkType := $(shell if [ -z ${NetworkType} ]; then echo "testnet"; else echo ${NetworkType}; fi)

BUILD_FLAGS = -ldflags "\
-X github.com/irisnet/irishub/version.GitCommit=${COMMIT_HASH} \
-X github.com/irisnet/irishub/types.InvariantLevel=${InvariantLevel} \
-X github.com/irisnet/irishub/types.NetworkType=${NetworkType}"

########################################
### Tools & dependencies

echo_bech32_prefix:
	@echo "\"source scripts/setProdEnv.sh\" to set compile environment variables for your product, or default values will be applied"
	@echo InvariantLevel=${InvariantLevel}
	@echo NetworkType=${NetworkType}

check_tools:
	cd scripts && $(MAKE) check_tools

check_dev_tools:
	cd scripts && $(MAKE) check_dev_tools

update_tools:
	cd scripts && $(MAKE) update_tools

update_dev_tools:
	cd scripts && $(MAKE) update_dev_tools

get_tools:
	cd scripts && $(MAKE) get_tools

get_dev_tools:
	cd scripts && $(MAKE) get_dev_tools

get_vendor_deps:
	@rm -rf vendor/
	@echo "--> Running dep ensure"
	@dep ensure -v

draw_deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i github.com/irisnet/irishub/cmd/iris -d 2 | dot -Tpng -o dependency-graph.png

########################################
### Generate swagger docs for irislcd
update_irislcd_swagger_docs:
	@statik -src=lite/swagger-ui -dest=lite -f

########################################
### Compile and Install
install: update_irislcd_swagger_docs echo_bech32_prefix
	go install $(BUILD_FLAGS) ./cmd/iris
	go install $(BUILD_FLAGS) ./cmd/iriscli
	go install $(BUILD_FLAGS) ./cmd/irislcd
	go install $(BUILD_FLAGS) ./cmd/iristool

build_linux: update_irislcd_swagger_docs echo_bech32_prefix
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o build/iris ./cmd/iris && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o build/iriscli ./cmd/iriscli && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o build/irislcd ./cmd/irislcd && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o build/iristool ./cmd/iristool

build_cur: update_irislcd_swagger_docs echo_bech32_prefix
	go build $(BUILD_FLAGS) -o build/iris ./cmd/iris  && \
	go build $(BUILD_FLAGS) -o build/iriscli ./cmd/iriscli && \
	go build $(BUILD_FLAGS) -o build/irislcd ./cmd/irislcd && \
	go build $(BUILD_FLAGS) -o build/iristool ./cmd/iristool


########################################
### Testing

test: test_unit test_cli test_lcd test_sim

test_sim: test_sim_modules test_sim_benchmark test_sim_iris_nondeterminism test_sim_iris_fast

test_unit:
	#@go test $(PACKAGES_NOSIMULATION)
	@go test $(PACKAGES_MODULES)
	@go test $(PACKAGES_TYPES)
	@go test $(PACKAGES_STORE)

test_cli:
	@go test  -timeout 20m -count 1 -p 4 tests/cli/*

test_lcd:
	@go test `go list github.com/irisnet/irishub/client/lcd`

test_sim_modules:
	@echo "Running individual module simulations..."
	@go test $(PACKAGES_SIMTEST)

test_sim_benchmark:
	@echo "Running benchmark test..."
	@go test ./app -run=none -bench=BenchmarkFullIrisSimulation -v -SimulationCommit=true -SimulationNumBlocks=100 -timeout 24h

test_sim_iris_nondeterminism:
	@echo "Running nondeterminism test..."
	@go test ./app -run TestAppStateDeterminism -v -SimulationEnabled=true -timeout 10m

test_sim_iris_fast:
	@echo "Running quick Iris simulation. This may take several minutes..."
	@go test ./app -run TestFullIrisSimulation -v -SimulationEnabled=true -SimulationNumBlocks=100 -timeout 24h

test_sim_iris_slow:
	@echo "Running full Iris simulation. This may take awhile!"
	@go test ./app -run TestFullIrisSimulation -v -SimulationEnabled=true -SimulationNumBlocks=1000 -SimulationVerbose=true -timeout 24h

testnet_init:
	@echo "Work well only when Bech32PrefixAccAddr equal faa"
	@if ! [ -f build/iris ]; then $(MAKE) build_linux ; fi
	@if ! [ -f build/nodecluster/node0/iris/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris testnet --v 4 --output-dir /home/nodecluster --chain-id irishub-test --starting-ip-address 192.168.10.2 ; fi
	@echo "To install jq command, please refer to this page: https://stedolan.github.io/jq/download/"
	@jq '.app_state.accounts+= [{"address": "faa1ljemm0yznz58qxxs8xyak7fashcfxf5lssn6jm", "coins": [ "1000000iris" ], "sequence_number": "0", "account_number": "0"}]' build/nodecluster/node0/iris/config/genesis.json > build/genesis_temp.json
	@sudo cp build/genesis_temp.json build/nodecluster/node0/iris/config/genesis.json
	@sudo cp build/genesis_temp.json build/nodecluster/node1/iris/config/genesis.json
	@sudo cp build/genesis_temp.json build/nodecluster/node2/iris/config/genesis.json
	@sudo cp build/genesis_temp.json build/nodecluster/node3/iris/config/genesis.json
	@rm build/genesis_temp.json
	@echo "Faucet address: faa1ljemm0yznz58qxxs8xyak7fashcfxf5lssn6jm"
	@echo "Faucet coin amount: 1000000iris"
	@echo "Faucet key seed: tube lonely pause spring gym veteran know want grid tired taxi such same mesh charge orient bracket ozone concert once good quick dry boss"

testnet_start:
	docker-compose up -d

testnet_stop:
	docker-compose down

testnet_clean:
	docker-compose down
	sudo rm -rf build/*

testnet_unsafe_reset:
	docker-compose down
	@docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris unsafe-reset-all --home=/home/nodecluster/node0/iris
	@docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris unsafe-reset-all --home=/home/nodecluster/node1/iris
	@docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris unsafe-reset-all --home=/home/nodecluster/node2/iris
	@docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris unsafe-reset-all --home=/home/nodecluster/node3/iris