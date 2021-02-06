DOCKER_HUB_TAG = "opsani/fiber-http:latest"
GITHUB_TAG = "ghcr.io/opsani/fiber-http:latest"

.PHONY: build
build:
	docker build -t $(DOCKER_HUB_TAG) -t $(GITHUB_TAG) .

.PHONY: run
run: build
	@mkdir -p ./build
	docker run -it --rm $(DOCKER_HUB_TAG)

.PHONY: push
push: build
	docker push $(DOCKER_HUB_TAG)
	docker push $(GITHUB_TAG)
