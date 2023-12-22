NAME=probum-users
VERSION=0.0.1

.PHONY: build
## build: Compile the packages and the API documentation.
build: docs
	@go build -o $(NAME)

.PHONY: run
## run: Build and Run in development mode.
run: build
	@./$(NAME) -e development

.PHONY: daemon
## daemon: Build and Run with CompileDaemon in development mode.
daemon: build
	@CompileDaemon --command="./$(NAME) -e development"

.PHONY: run-prod
## run-prod: Build and Run in production mode.
run-prod: build
	@./$(NAME) -e production

.PHONY: clean
## clean: Clean project, the generated documentation and previous builds.
clean:
	-@rm docs/swagger.*
	-@rm docs/docs.go
	@rm -f $(NAME)

.PHONY: docs
## docs: Generate Swagger documentation
docs:
	@swag init -g docs/docgen.go -o ./docs --parseDependency true

.PHONY: deps
## deps: Download modules
deps:
	@go mod download

.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo