PACKAGE_PATH = ./...
OUTPUT_PATH = cmd/

.PHONY: build
build:
	go  build -o ${OUTPUT_PATH} ${PACKAGE_PATH}
