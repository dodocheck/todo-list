gen:
	protoc --proto_path=proto \
	proto/tasks.proto proto/service.proto \
	--go_out=pb \
	--go_opt=paths=source_relative \
	--go-grpc_out=pb \
	--go-grpc_opt=paths=source_relative

clean:
	rm -rf pb/*
