AGENT_DIR=./cmd/agent
SERVER_DIR=./cmd/server

build-agent: $(AGENT_DIR)/main.go
	go build -buildvcs=false -o $(AGENT_DIR)/agent $(AGENT_DIR)

build-server: $(SERVER_DIR)/main.go
	go build -buildvcs=false -o $(SERVER_DIR)/server $(SERVER_DIR) 

build: build-agent build-server

agent: $(AGENT_DIR)/agent
	$(AGENT_DIR)/agent -a 0.0.0.0:8080 -p 5 -r 5

server: $(SERVER_DIR)/server
	$(SERVER_DIR)/server -a :8080
