## Run Test Suite inside docker
tests:
	@docker-compose -f docker-compose.yml up --build --abort-on-container-exit

## Run Integration Test
## Note: This command is intended to be executed within docker env
integration-tests:
	go test -v ./...