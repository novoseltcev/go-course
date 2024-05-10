AGENT_DIR=./cmd/agent
SERVER_DIR=./cmd/server

generate:
	go generate ./...

build-agent: generate $(AGENT_DIR)/main.go
	go build -buildvcs=false -o $(AGENT_DIR)/agent $(AGENT_DIR)

build-server: generate $(SERVER_DIR)/main.go
	go build -buildvcs=false -o $(SERVER_DIR)/server $(SERVER_DIR) 

build: build-agent build-server

agent: $(AGENT_DIR)/agent
	$(AGENT_DIR)/agent -a 0.0.0.0:8080 -p 2 -r 5 -k secret-key 

server: $(SERVER_DIR)/server
	KEY=secret-key DATABASE_DSN="postgres://postgres:postgres@0.0.0.0:5432/praktikum?sslmode=disable" RESTORE=true STORE_INTERVAL=2 FILE_STORAGE_PATH=$(SERVER_DIR)/backup.json $(SERVER_DIR)/server -a :8080

up:
	docker-compose up -d --build

down:
	docker-compose down

psql:
	PGPASSWORD=postgres psql -U postgres -h 0.0.0.0 -p 5432 -d test
