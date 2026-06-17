# Ma'at — one-command verification.
# `make test` runs everything CI runs, locally.

.PHONY: test sim sim-real chain-test demo fetch calibrate clean

test: chain-test sim ## run all verification (Go core + economic gate)
	@echo "ALL CHECKS PASSED"

chain-test: ## Go: vet + unit/integration tests for the entire project
	go vet ./... && go test ./... -count=1

sim: ## Python: economic stop-gate (synthetic)
	python3 simulation/maat_sim.py --paths 1000 --seed 42

sim-real: ## Python: economic stop-gate driven by real ETH data
	cd simulation && python3 maat_sim.py --paths 1000 --seed 42 \
		--price-csv data/eth_hourly.csv --annual-vol 0.4554 --jump-size 0.0291 --crash-size 0.38

demo: ## run the end-to-end chain demo (oracle+market+treasury+bridge)
	cd chain && go run ./cmd/demo

fetch: ## refresh real ETH price data (CoinGecko)
	python3 simulation/fetch_data.py

calibrate: ## estimate sim parameters from real data
	cd simulation && python3 calibrate.py

clean:
	rm -rf simulation/charts simulation/calibrated_params.json

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN{FS=":.*?## "}{printf "  %-12s %s\n", $$1, $$2}'
