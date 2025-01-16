main: generate fix lint test cover build up migrate server

generate:
	go generate ./...

AGENT_DIR=./cmd/agent

build-agent: generate $(AGENT_DIR)/main.go
	go build -buildvcs=false -ldflags "-X main.buildVersion=v1.0.0 -X main.buildDate=`date -u +%Y-%m-%d` -X main.buildCommit=`git rev-parse HEAD`" -o $(AGENT_DIR)/agent $(AGENT_DIR)

SERVER_DIR=./cmd/server

build-server: generate $(SERVER_DIR)/main.go
	go build -buildvcs=false -ldflags "-X main.buildVersion=v1.0.0 -X main.buildDate=`date -u +%Y-%m-%d` -X main.buildCommit=`git rev-parse HEAD`" -o $(SERVER_DIR)/server $(SERVER_DIR) 

STATICLINT_DIR=./cmd/staticlint

build-staticlint: $(STATICLINT_DIR)/main.go
	go build -buildvcs=false -o $(STATICLINT_DIR)/staticlint $(STATICLINT_DIR)

build: build-agent build-server build-staticlint


.PHONY: staticlint
staticlint: build-staticlint $(STATICLINT_DIR)/staticlint
	$(STATICLINT_DIR)/staticlint ./...

.PHONY: agent
agent: $(AGENT_DIR)/agent
	$(AGENT_DIR)/agent -c ./config/agent.json

DATABASE_URI=postgresql://postgres:postgres@0.0.0.0:5432/praktikum?sslmode=disable
.PHONY: server
server: $(SERVER_DIR)/server
	$(SERVER_DIR)/server -d $(DATABASE_URI) --config ./config/server.json

migrate:
	migrate -source file://migrations -database $(DATABASE_URI) up


up:
	docker-compose up -d --build

down:
	docker-compose down

.PHONY: psql
psql:
	psql $(DATABASE_URI)


fix:
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

