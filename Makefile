PROTOCOLS := $(shell find protocol/msg -name "*.proto" -print)

gen:
	protoc -I=protocol/msg --go_out=gen --go_opt=paths=import $(PROTOCOLS)

run:
	go run src/main.go