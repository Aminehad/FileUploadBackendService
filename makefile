UP_SERVICES?=
OPENAPI_FILE?=internal/app/upload-service/api/openapi.yml


up: ## Build and Run all the docker containers of the project (upload-service / other services...)
	@docker compose up --remove-orphans --build -d --wait ${UP_SERVICES}
	@docker compose logs -f ${UP_SERVICES}

stop:  ## stop all the project docker containers
	@docker compose stop

down:  ## Delete all the project docker containers
	@docker compose down

lint:
	@golangci-lint run ./...

openapi:
	npx @redocly/cli preview-docs --port 8083 ${OPENAPI_FILE}