# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

start: ## Run docker compose up -d
ifdef name
	docker compose up -d $(name)
else
	docker compose up -d
endif
.PHONY: start

stop: ## Run docker compose down --remove-orphans
ifdef name
	docker compose down --remove-orphans $(name)
else
	docker compose down --remove-orphans
endif
.PHONY: stop

restart: ## Run docker compose down --remove-orphans && Run docker compose up -d
ifdef name
	$(MAKE) stop name=$(name)
	$(MAKE) start name=$(name)
else
	$(MAKE) stop
	$(MAKE) start
endif
.PHONY: restart

build: ## Run docker compose build
ifdef file
	docker compose -f $(file) build
else
	docker compose build
endif
.PHONY: build

swag-generate: ## Run swag init -g ./internal/handler/http/v1/v1.go
	docker compose exec ten_tarot_server swag init -g ./internal/handler/http/v1/v1.go
.PHONY: swag

sh: ## Run docker compose exec $(name) sh
ifdef name
	docker compose exec $(name) sh
else
	@echo "name is required"
	exit 1
endif
.PHONY: sh

logs:
ifdef name
	docker compose logs -f $(name)
else
	docker compose logs -f
endif
.PHONY: logs