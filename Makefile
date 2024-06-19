CMD = ./cmd
TEST = ./test
BIN = ./bin/liteapi
DOCS_DIR = ./cmd/api/docs

install:
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest

clean:
	rm -rf ./bin

docs:
	swag init --dir ./cmd --output ${DOCS_DIR}

dev:
	go run ${CMD}

build:
	go build -o ${BIN} ${CMD}

start: clean build
	${BIN}
