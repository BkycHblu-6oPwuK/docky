# Makefile

BINARY_NAME=docky
BUILD_DIR=bin
APP=cmd/docky/main.go

build:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) $(APP)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(APP)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(APP)

run: 
	go run $(APP) $(filter-out $@,$(MAKECMDGOALS))

clean:
	rm -rf $(BUILD_DIR)

%:
	@: