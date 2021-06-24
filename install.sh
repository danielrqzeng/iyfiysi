#statik -src=template; go build -o iyfiysi iyfiysi.go
go mod download

go install github.com/rakyll/statik
go install github.com/golang/protobuf/protoc-gen-go
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

BUILT_AT=$(date "+%Y-%m-%d %H:%M:%S") #2020-03-30 10:08:54
COMMIT_TAG=$(git rev-parse HEAD) #1c7caa847ce196f0668e01794d3cd773944f3127
if [ ${#COMMIT_TAG} -eq 40 ];then
    COMMIT_TAG=${COMMIT_TAG:0:8}
else
    COMMIT_TAG="unknow"
fi

statik -src=template -f
cd cmd/iyfiysi;go install -ldflags "-X 'iyfiysi/internal/iyfiysi.commit=$COMMIT_TAG' -X 'iyfiysi/internal/iyfiysi.date=$BUILT_AT'" .
cd -

cd cmd/protoc-gen-iyfiysi;go install -ldflags "-X 'iyfiysi/internal/iyfiysi.commit=$COMMIT_TAG' -X 'iyfiysi/internal/iyfiysi.date=$BUILT_AT'" .
cd -