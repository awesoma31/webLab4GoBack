gen:
	@protoc --go_out=. --go_opt=paths=source_relative \
         --go-grpc_out=. --go-grpc_opt=paths=source_relative \
         api/oms.proto  \
         api/model.proto \
         api/pointsService.proto