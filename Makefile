init: tidy proto

proto:
	cd api/protobuf/ && make gen

tidy:
	go mod tidy

