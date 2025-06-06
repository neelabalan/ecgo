BINARY_NAME = ecgo
IMAGE_NAME = ecgo
VERSION ?= 0.1.0
PLATFORMS = linux/amd64,linux/arm64,linux/arm/v7

.PHONY: build test run clean docker-build docker-push docker-build-local docker-test-local docker-test-multi-arch

build:
	mkdir -p bin
	go build -o bin/${BINARY_NAME} -ldflags "-X main.version=${VERSION}" ecgo.go

test:
	go test -v ./...

test-cov:
	mkdir -p coverage
	go test -v -cover -coverprofile=coverage/coverage.out ./...

run: build
	./bin/${BINARY_NAME}

clean:
	go clean
	rm -rf bin
	rm -f ${BINARY_NAME}

docker-build:
	# Build multi-architecture images and save them locally
	docker buildx build \
		--platform ${PLATFORMS} \
		-t ${IMAGE_NAME}:${VERSION} \
		-t ${IMAGE_NAME}:latest \
		--build-arg VERSION=${VERSION} \
		.

docker-push:
	# Build and push multi-architecture images to a registry
	docker buildx build \
		--platform ${PLATFORMS} \
		-t ${IMAGE_NAME}:${VERSION} \
		-t ${IMAGE_NAME}:latest \
		--build-arg VERSION=${VERSION} \
		--push .

docker-build-local:
	# Standard local Docker build for the host architecture
	docker build \
		-t ${IMAGE_NAME}:${VERSION}-local \
		-t ${IMAGE_NAME}:latest-local \
		--build-arg VERSION=${VERSION} \
		.

docker-test-local:
	# Build a test image for the host architecture and run tests
	docker build \
		-t ${IMAGE_NAME}-test:local \
		-f Dockerfile.test .
	docker create --name ${IMAGE_NAME}-test-container ${IMAGE_NAME}-test:local
	mkdir -p coverage
	docker cp ${IMAGE_NAME}-test-container:/coverage.out coverage/coverage.out
	docker rm ${IMAGE_NAME}-test-container
	@echo "Coverage report generated: coverage/coverage.out"
	@echo "Run 'go tool cover -html=coverage/coverage.out' to view the report."

docker-test-multi-arch:
	# Iterate through the defined PLATFORMS for multi-arch testing
	@for platform in $(subst ,,$(PLATFORMS)); do \
		echo "--- Testing on platform: $$platform ---"; \
		# Build the test image for the specific platform using buildx
		docker buildx build \
			--platform $$platform \
			-t ${IMAGE_NAME}-test:$$platform \
			-f Dockerfile.test \
			--load . ; \
		# Create and run a temporary container for this platform
		docker create --name ${IMAGE_NAME}-test-container-$$platform ${IMAGE_NAME}-test:$$platform; \
		# Create a platform-specific coverage directory
		mkdir -p coverage/$$platform; \
		# Copy the coverage report
		docker cp ${IMAGE_NAME}-test-container-$$platform:/coverage.out coverage/$$platform/coverage.out; \
		# Clean up the container
		docker rm ${IMAGE_NAME}-test-container-$$platform; \
		echo "Coverage report for $$platform generated: coverage/$$platform/coverage.out"; \
		echo ""; \
	done
	@echo "All multi-architecture local tests completed."
	@echo "View reports: 'go tool cover -html=coverage/<platform>/coverage.out'"