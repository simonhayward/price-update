.PHONY: test
test: ## run local tests
	go test -race ./...

.PHONY: trace
trace: ## run trace tool
	go test -trace trace.out ./...priceupdate

.PHONY: build
build: ## build binary
	go build

.PHONY: docker-build
docker-build: ## docker build
	docker build -t price-update .

.PHONY: docker-run
docker-run: docker-build ## docker run in container
	docker run --env INPUT=${INPUT} --env OUTPUT=${OUTPUT} \
	--env TOKEN=${TOKEN} --env API=${API} \
	--rm price-update

.PHONY: deploy
deploy: ## deploy out as GCP function
	gcloud functions deploy priceupdate \
	--set-env-vars INPUT=${INPUT},OUTPUT=${OUTPUT},TOKEN=${TOKEN},API=${API} \
	--entry-point RunUpdate --runtime go113 --trigger-http --memory=128MB --region=europe-west2 \
	--source=./priceupdate --timeout=20s

.PHONY: help
help:  ## help command
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
