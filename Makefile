## Run Test Suite inside docker
unit-tests:
	go tests -v -short ./...

## Run Integration Test
## Note: This command is intended to be executed within docker env
integration-tests:
	go test -run Integration ./...

## Run Integration Test inside docker
integration-tests-docker:
	docker-compose up --build