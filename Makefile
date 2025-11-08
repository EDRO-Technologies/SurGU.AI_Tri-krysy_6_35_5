.PHONY: m-create
m-create:
ifndef NAME
	$(error NAME is not set. Usage: make m-create NAME=add_users_table)
endif
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	touch migrations/$${timestamp}_$(NAME).up.sql; \
    touch migrations/$${timestamp}_$(NAME).down.sql;

.PHONY: fmt
fmt:
	docker run --rm -v "$(CURDIR):/app" -w /app golangci/golangci-lint:latest golangci-lint run --fix