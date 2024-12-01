main: generate format lint build up migrate server

generate:
	go generate ./...

AGENT_DIR=./cmd/agent

build-agent: generate $(AGENT_DIR)/main.go
	go build -buildvcs=false -o $(AGENT_DIR)/agent $(AGENT_DIR)

SERVER_DIR=./cmd/server

build-server: generate $(SERVER_DIR)/main.go
	go build -buildvcs=false -o $(SERVER_DIR)/server $(SERVER_DIR) 

STATICLINT_DIR=./cmd/staticlint

build-staticlint: $(STATICLINT_DIR)/main.go
	go build -buildvcs=false -o $(STATICLINT_DIR)/staticlint $(STATICLINT_DIR)

build: build-agent build-server build-staticlint


.PHONY: staticlint
staticlint: build-staticlint $(STATICLINT_DIR)/staticlint
	$(STATICLINT_DIR)/staticlint ./...

.PHONY: agent
agent: $(AGENT_DIR)/agent
	$(AGENT_DIR)/agent -a 0.0.0.0:8080 -p 2 -r 5

DATABASE_URI=postgresql://postgres:postgres@0.0.0.0:5432/praktikum?sslmode=disable
.PHONY: server
server: $(SERVER_DIR)/server
	DATABASE_DSN=$(DATABASE_URI) RESTORE=true STORE_INTERVAL=2 FILE_STORAGE_PATH=$(SERVER_DIR)/backup.json $(SERVER_DIR)/server -a :8080

migrate:
	migrate -source file://migrations -database $(DATABASE_URI) up


up:
	docker-compose up -d --build

down:
	docker-compose down

.PHONY: psql
psql:
	psql $(DATABASE_URI)


format:
	golangci-lint run --fix

lint: build-staticlint
	$(STATICLINT_DIR)/staticlint ./... & golangci-lint run

COVERAGE_PROFILE=reports/profile.cov
test:
	go clean -testcache
	go test -coverprofile=$(COVERAGE_PROFILE) -bench=. -benchmem ./...
	grep -v -E -f .covignore $(COVERAGE_PROFILE) > $(COVERAGE_PROFILE).filtered && mv $(COVERAGE_PROFILE).filtered $(COVERAGE_PROFILE)

cover:
	go tool cover -func=$(COVERAGE_PROFILE) -o reports/coverage.out

cover-html:
	go tool cover -html=$(COVERAGE_PROFILE) -o reports/coverage.html

docs:
	pkgsite -http=:8080

