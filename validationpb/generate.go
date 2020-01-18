package validationpb

//go:generate protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf --gofast_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,paths=source_relative:. validation.proto
