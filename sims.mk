#!/usr/bin/make -f

########################################
### Simulations

SIMAPP = github.com/bacd/app

sim-bac-nondeterminism:
	@echo "Running nondeterminism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -SimulationEnabled=true -v -timeout 10m

sim-bac-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.bacd/config/genesis.json will be used."
	@go test -mod=readonly github.com/bac/app -run TestFullBacSimulation -SimulationGenesis=${HOME}/.bacd/config/genesis.json \
		-SimulationEnabled=true -SimulationNumBlocks=100 -SimulationBlockSize=200 -SimulationCommit=true -SimulationSeed=99 -SimulationPeriod=5 -v -timeout 24h

sim-bac-fast:
	@echo "Running quick Bac simulation. This may take several minutes..."
	@go test -mod=readonly github.com/bac/app -run TestFullBacSimulation -SimulationEnabled=true -SimulationNumBlocks=100 -SimulationBlockSize=200 -SimulationCommit=true -SimulationSeed=99 -SimulationPeriod=5 -v -timeout 24h

sim-bac-import-export: runsim
	@echo "Running Bac import/export simulation. This may take several minutes..."
	$(GOPATH)/bin/runsim 25 5 TestBacImportExport

sim-bac-simulation-after-import: runsim
	@echo "Running Bac simulation-after-import. This may take several minutes..."
	$(GOPATH)/bin/runsim 25 5 TestBacSimulationAfterImport

sim-bac-custom-genesis-multi-seed: runsim
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.bacd/config/genesis.json will be used."
	$(GOPATH)/bin/runsim -g ${HOME}/.bacd/config/genesis.json 400 5 TestFullBacSimulation

sim-bac-multi-seed: runsim
	@echo "Running multi-seed Bac simulation. This may take awhile!"
	$(GOPATH)/bin/runsim 400 5 TestFullBacSimulation

sim-benchmark-invariants:
	@echo "Running simulation invariant benchmarks..."
	@go test -mod=readonly github.com/bac/app -benchmem -bench=BenchmarkInvariants -run=^$ \
	-SimulationEnabled=true -SimulationNumBlocks=1000 -SimulationBlockSize=200 \
	-SimulationCommit=true -SimulationSeed=57 -v -timeout 24h

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true
sim-bac-benchmark:
	@echo "Running Bac benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ github.com/bac/app -bench ^BenchmarkFullBacSimulation$$  \
		-SimulationEnabled=true -SimulationNumBlocks=$(SIM_NUM_BLOCKS) -SimulationBlockSize=$(SIM_BLOCK_SIZE) -SimulationCommit=$(SIM_COMMIT) -timeout 24h

sim-bac-profile:
	@echo "Running Bac benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ github.com/bac/app -bench ^BenchmarkFullBacSimulation$$ \
		-SimulationEnabled=true -SimulationNumBlocks=$(SIM_NUM_BLOCKS) -SimulationBlockSize=$(SIM_BLOCK_SIZE) -SimulationCommit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out


.PHONY: runsim sim-bac-nondeterminism sim-bac-custom-genesis-fast sim-bac-fast sim-bac-import-export \
	sim-bac-simulation-after-import sim-bac-custom-genesis-multi-seed sim-bac-multi-seed \
	sim-benchmark-invariants sim-bac-benchmark sim-bac-profile
