BIN_NAME="guff"

build:
	go build -o ${BIN_NAME}

clean:
	go clean
	rm ./${BIN_NAME}
	rm /usr/bin/${BIN_NAME}

install: build
	cp ${BIN_NAME} /usr/bin/

uninstall:
	rm /usr/bin/${BIN_NAME}

run: build
	./${BIN_NAME}
