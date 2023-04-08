USER_BINARY=userSvc
## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_user_svc
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"
## build_user_svc: builds the user binary as a linux executable
build_user_svc:
	@echo "Building user binary..."
	cd ../user-svc && env GOOS=linux CGO_ENABLED=0 go build -o ${USER_BINARY} ./cmd/api
	@echo "Done!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"
