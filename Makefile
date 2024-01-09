VERSION := 3.12.11
TAG := $(shell echo $(VERSION) | cut -d '.' -f -2)

build-cluster-bootstrap:
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o rabbitmq/dist/cluster-bootstrap ./cluster-bootstrap

build-docker:
	docker build --build-arg RABBITMQ_VERSION=$(VERSION) \
		-t upfluence/rabbitmq:$(TAG) ./rabbitmq

build: build-cluster-bootstrap build-docker
