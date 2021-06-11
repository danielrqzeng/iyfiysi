#statik -src=template; go build -o iyfiysi iyfiysi.go
go mod download

go install github.com/rakyll/statik
go install github.com/golang/protobuf/protoc-gen-go
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

statik -src=template
cd cmd/iyfiysi;go install .
cd -
cd cmd/protoc-gen-iyfiysi;go install .
cd -