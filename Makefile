main: generate build up migrate server

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

test:
	go test -v -coverprofile=profiles/profile.cov -bench=. -benchmem ./...

cover: test
	go tool cover -func=profiles/profile.cov

cover-html: test
	go tool cover -html=profiles/profile.cov
