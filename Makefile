main: generate format lint build up migrate server

format:
	golangci-lint run --fix

lint:
	golangci-lint run

fast-lint:
	golangci-lint run --fast	

generate:
	go generate ./...

AGENT_DIR=./cmd/agent

build-agent: generate $(AGENT_DIR)/main.go
	go build -buildvcs=false -o $(AGENT_DIR)/agent $(AGENT_DIR)

SERVER_DIR=./cmd/server

build-server: generate $(SERVER_DIR)/main.go
	go build -buildvcs=false -o $(SERVER_DIR)/server $(SERVER_DIR) 

build: build-agent build-server

agent: $(AGENT_DIR)/agent
	$(AGENT_DIR)/agent -a 0.0.0.0:8080 -p 2 -r 5 -k secret-key 

DATABASE_URI=postgresql://postgres:postgres@0.0.0.0:5432/praktikum?sslmode=disable

server: $(SERVER_DIR)/server
	KEY=secret-key DATABASE_DSN=$(DATABASE_URI) RESTORE=true STORE_INTERVAL=2 FILE_STORAGE_PATH=$(SERVER_DIR)/backup.json $(SERVER_DIR)/server -a :8080

up:
	docker-compose up -d --build

down:
	docker-compose down

.PHONY: psql
psql:
	psql $(DATABASE_URI)

migrate:
	migrate -source file://migrations -database $(DATABASE_URI) up

COVERAGE_PROFILE=reports/profile.cov
test:
	go test -v -coverprofile=$(COVERAGE_PROFILE) -bench=. -benchmem ./...
	grep -v -E -f .covignore $(COVERAGE_PROFILE) > $(COVERAGE_PROFILE).filtered && mv $(COVERAGE_PROFILE).filtered $(COVERAGE_PROFILE)

cover:
	go tool cover -func=$(COVERAGE_PROFILE)

cover-html:
	go tool cover -html=$(COVERAGE_PROFILE)

docs:
	pkgsite -http=:8080
