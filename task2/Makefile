container_runtime := $(shell which podman || which docker)

$(info using ${container_runtime})

up: down
	${container_runtime} compose up --build -d

down:
	${container_runtime} compose down

run-tests: 
	${container_runtime} run --rm --network=host tests:latest

test:
	make down
	make up
	make run-tests
	make down
	@echo "test finished"

lint:
	make -C petname lint

proto:
	make -C petname protobuf
