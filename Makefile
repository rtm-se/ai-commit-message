BINARY_NAME = auto-commit

clean:
	@rm -rfd ./bin

build:
	go mod tidy
	mkdir "bin"
	go build -o bin/${BINARY_NAME} cmd/ai-commit/main.go

rebuild:
	make clean
	make build