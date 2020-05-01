DOCKER_ACCOUNT = barcodepro
SITENAME = weaponry
APPNAME = go-service-template

COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS = -a -installsuffix cgo -ldflags "-X main.appName=${APPNAME} -X main.gitCommit=${COMMIT} -X main.gitBranch=${BRANCH}"

.PHONY: help \
		clean lint test race \
		build migrate docker-build docker-push deploy

.DEFAULT_GOAL := help

help: ## Display this help screen
	@echo "Makefile available targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  * \033[36m%-15s\033[0m %s\n", $$1, $$2}'

clean: ## Clean
	rm -f bin/${APPNAME}

dep: ## Get the dependencies
	go mod download

lint: ## Lint the source files
	golangci-lint run --timeout 5m -E golint -e '(method|func) [a-zA-Z]+ should be [a-zA-Z]+'

test: dep ## Run unittests
	go test -short -timeout 300s -p 1 ./...

race: dep ## Run data race detector
	go test -race -short -timeout 300s -p 1 ./...

#coverage: ## Generate global code coverage report
#  ./tools/coverage.sh;
#
#coverhtml: ## Generate global code coverage report in HTML
#  ./tools/coverage.sh html;

build: dep ## Build executable
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${APPNAME} ./service/cmd

#migrate: ## Run migrations
#	migrate -database ${URI} -path migrations/ up

#docker-build: ## Build docker image
#	docker build -t ${DOCKER_ACCOUNT}/${SITENAME}-${APPNAME}:${COMMIT} .
#	docker image prune --force --filter label=stage=intermediate
#	docker tag ${DOCKER_ACCOUNT}/${SITENAME}-${APPNAME}:${COMMIT} ${DOCKER_ACCOUNT}/${SITENAME}-${APPNAME}:latest
#
#docker-push: ## Push docker image to registry
#	docker push ${DOCKER_ACCOUNT}/${SITENAME}-${APPNAME}:${COMMIT}
#	docker push ${DOCKER_ACCOUNT}/${SITENAME}-${APPNAME}:latest

#deploy: ## Deploy application
#	ansible-playbook --vault-password-file=${ANSIBLE_VAULT_PASSWORD_FILE} deployment/ansible/deploy.yml -e env=${ENV}
