# Variables
GO_CMD=go
PROTOC_CMD=protoc

# Directories
GO_DIR=./go
PROTO_DIR=./proto
GO_OUT_DIR=$(GO_DIR)/generated

# Proto files
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

# Default target
.PHONY: all
all: generate-go

# Generate Go code from Proto files
.PHONY: generate-go
generate-go: $(PROTO_FILES)
	@echo "Generating Go code from Proto files..."
	@mkdir -p $(GO_OUT_DIR)
	@$(PROTOC_CMD) --go_out=$(GO_OUT_DIR) --go-grpc_out=$(GO_OUT_DIR) --gorm_out=$(GO_OUT_DIR) --proto_path=$(PROTO_DIR) $(PROTO_FILES)
	@cd $(GO_OUT_DIR)/prototiger && go mod init prototiger
	@echo "$(PROTO_FILES) GENERATED as GOLANG IN $(GO_OUT_DIR)"

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning Go build files..."
	@cd $(GO_DIR) && $(GO_CMD) clean
	@rm -rf $(GO_OUT_DIR)

re: clean all