gen: gen-user gen-group gen-aiProxy

gen-user:
	protoc --go_out=../  --go-grpc_out=../  ./user/*.proto

gen-group:
	protoc --go_out=../ --go-grpc_out=../ ./group/*.proto --experimental_allow_proto3_optional

gen-aiProxy:
	protoc --go_out=../ --go-grpc_out=../ ./ai-proxy/*.proto