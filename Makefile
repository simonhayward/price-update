.PHONY: test
test: ## run local tests
	go test -race ./... -v

.PHONY: build
build: ## build binary
	go build

.PHONY: docker-build
docker-build: ## docker build
	docker build -t price-update .

.PHONY: docker-run
docker-run: docker-build ## docker run in container
	docker run --env SOURCES=${SOURCES} --env STORE=${STORE} \
	--env STORE_TOKEN=${STORE_TOKEN} --env STORE_UPDATE=${STORE_UPDATE} \
	--rm price-update

.PHONY: deploy
deploy: ## deploy out as GCP function
	gcloud functions deploy priceupdate \
	--set-env-vars SOURCES=${SOURCES},STORE=${STORE},STORE_TOKEN=${STORE_TOKEN},STORE_UPDATE=${STORE_UPDATE} \
	--entry-point RunUpdate --runtime go113 --trigger-http --memory=128MB --region=europe-west2 \
	--source=./priceupdate --timeout=10s

.PHONY: help
help:  ## help command
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
