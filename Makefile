IMAGE_NAME = "opsani/fiber-http:latest"

.PHONY: build
build:
	docker build -t $(IMAGE_NAME) .

.PHONY: run
run: build
	@mkdir -p ./build
	docker run -it --rm $(IMAGE_NAME)

.PHONY: push
push: build
	docker push $(IMAGE_NAME)
