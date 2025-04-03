BINARY_NAME = auto-commit

clean:
	@rm -rfd ./bin

build:
	mkdir "bin"
	go build -o bin/${BINARY_NAME} cmd/ai-commit/main.go


