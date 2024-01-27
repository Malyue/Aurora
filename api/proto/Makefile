gen: gen-user gen-group

gen-user:
	protoc --go_out=../  --go-grpc_out=../  ./user/*.proto

gen-group:
	protoc --go_out=../ --go-grpc_out=../ ./group/*.proto --experimental_allow_proto3_optional