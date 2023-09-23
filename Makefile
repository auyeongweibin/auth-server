generate_grpc_code:
	protoc \
	--go_out=registration \
	--go_opt=paths=source_relative \
	--go-grpc_out=registration \
	--go-grpc_opt=paths=source_relative \
	registration.proto && \
	protoc \
	--go_out=otp \
	--go_opt=paths=source_relative \
	--go-grpc_out=otp \
	--go-grpc_opt=paths=source_relative \
	otp.proto