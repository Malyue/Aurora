PROTO_PATH ?= api/proto

init: tidy proto

run:

build:

proto:
	cd "${PROTO_PATH}" && make gen

tidy:
	go mod tidy

