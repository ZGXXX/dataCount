module github.com/grpc-demo/datacount/client

go 1.15

require (
	github.com/grpc-demo/datacount/protoc v1.0.0
	google.golang.org/grpc v1.39.0
)

replace github.com/grpc-demo/datacount/protoc => ../protoc
