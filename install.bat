echo off
echo ==========================================
echo install iyfiysi and its dependency
echo on

rem BUILT_AT=2020-03-30 10:08:54
set BUILT_AT=%date:~0,4%-%date:~5,2%-%date:~8,2% %time%
set COMMIT_TAG=unknow
for /F %%i in ('git rev-parse HEAD') do ( set COMMIT_TAG=%%i)

@rem run me to install iyfiysi and its dependency
go mod download

@rem install statik v0.1.7
go install github.com/rakyll/statik

@rem install protoc-gen-go v1.3.2
go install github.com/golang/protobuf/protoc-gen-go

@rem install protoc-gen-grpc-gateway v1.13.0
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

@rem install protoc-gen-grpc-gateway v1.13.0
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

@rem install protoc-gen-swagger v1.13.0
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

@rem install iyfiysi & protoc-gen-iyfiysi
%GOPATH%\bin\statik -src=template -f
cd cmd\iyfiysi
go install -ldflags "-X 'iyfiysi/internal/iyfiysi.commit=%COMMIT_TAG%' -X 'iyfiysi/internal/iyfiysi.date=%BUILT_AT%'" .

cd ..\protoc-gen-iyfiysi
go install -ldflags "-X 'iyfiysi/internal/iyfiysi.commit=%COMMIT_TAG%' -X 'iyfiysi/internal/iyfiysi.date=%BUILT_AT%'" .

cd ..
cd ..
@echo off
echo ==========================================