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

```yaml
flags_over_config: true

ignore_patterns:
  - message: "updated proto"
    patterns:
      - ".*\\.pb\\.go$"
      - ".*\\.pb\\.gw\\.go$"
      - ".*\\.swagger\\.json"
      - ".*\\.pb\\.framework\\.go"
  ```