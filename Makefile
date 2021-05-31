.PHONY: help

help: ## help command for available tasks
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help


build: ## build the container
	docker build -t cftc-cot .

build-nc: ## build the container w/o a cache
	docker build --no-cache -t cftc-cot .

run: ## run the container with default parameters
	docker run --net=host -v $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/src/config/:/config cftc-cot

up: build run ## build the container and boot

image:
	docker tag cftc-cot docker.pkg.github.com/adityaxdiwakar/cftc-cot/cftc-cot:${TRAVIS_TAG}
	docker tag cftc-cot docker.pkg.github.com/adityaxdiwakar/cftc-cot/cftc-cot:latest
	
push-image:
	docker push docker.pkg.github.com/adityaxdiwakar/cftc-cot/cftc-cot:${TRAVIS_TAG}
	docker push docker.pkg.github.com/adityaxdiwakar/cftc-cot/cftc-cot:latest
