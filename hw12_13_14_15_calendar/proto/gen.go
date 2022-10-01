package proto

//go:generate protoc --go_out=./event --go-grpc_out=./event ./event/event.proto
//go:generate protoc -I . --grpc-gateway_out ./ --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true ./event/event.proto
