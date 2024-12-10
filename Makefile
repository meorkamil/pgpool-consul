BIN_NAME=pgpool-consul
CMD_DIR=cmd
VERSION=v0.0.0
BUILD_DIR=build
CONFIG_PATH=./config/config.yml

debug:
	cd ${CMD_DIR}/${BIN_NAME} && go run main.go

build:
	CGO_ENABLED=0 go build -C ${CMD_DIR}/${BIN_NAME}  -o ../../${BUILD_DIR}/${BIN_NAME}
	tar -czf ${BUILD_DIR}/${BIN_NAME}-${VERSION}.tar.gz ${BUILD_DIR}/${BIN_NAME}

run: build
	./${BUILD_DIR}/${BIN_NAME} --config ${CONFIG_PATH}

clean-test:
	@go clean -testcache

test:
	go test ./...

clean:
	rm -rf ${BUILD_DIR}
