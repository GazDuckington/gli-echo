PACKAGE_PATH = ./...
OUTPUT_PATH = tmp/

.PHONY: build
build:
	go  build -o ${OUTPUT_PATH} ${PACKAGE_PATH}
