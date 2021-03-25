#statik -src=template; go build -o iyfiysi iyfiysi.go
statik -src=template
cd cmd/iyfiysi;go install .
cd -
cd cmd/protoc-gen-iyfiysi;go install .
cd -