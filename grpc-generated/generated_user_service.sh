export GOPATH=~/go # set your own go path
export PATH=$PATH:$GOPATH/bin
protoc --go_out=./go_package_files --go-grpc_out=./go_package_files proto_files/user_service.proto
