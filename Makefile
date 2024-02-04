
run-handler:
	@go run eventhandler

run-processor:
	@go run eventprocessor
	
run:
	@docker-compose up --pull .
