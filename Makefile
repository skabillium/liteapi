CMD = ./cmd
TEST = ./test
BIN = ./bin/liteapi
DOCS_DIR = ./cmd/api/docs

install:
	go mod download

clean:
	rm -rf ./bin

docs:
	swag init --dir ./cmd --output ${DOCS_DIR}

tests:
	go test ${CMD}/db
	go test ${CMD}/resp
	go test ${CMD}

dev:
	go run ${CMD}

build:
	go build -o ${BIN} ${CMD}

start: clean build
	${BIN}
