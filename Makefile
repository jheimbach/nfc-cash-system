SRC_DIR=server/api
DST_DIR=$(SRC_DIR)

generateGrpc:
	protoc -I=$(SRC_DIR) \
		-I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.12.1/third_party/googleapis \
    	-I=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.12.1 \
		--go_out=plugins=grpc:$(DST_DIR) \
		$(SRC_DIR)/*.proto
generateGateway:
	protoc -I=$(SRC_DIR) \
		-I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.12.1/third_party/googleapis \
    	-I=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.12.1 \
		--grpc-gateway_out=logtostderr=true:$(DST_DIR) \
		$(SRC_DIR)/*.proto

generateSwagger:
	protoc -I=$(SRC_DIR) \
    		-I=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.12.1/third_party/googleapis \
    		-I=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.12.1 \
    		--swagger_out=logtostderr=true,allow_merge=true:$(DST_DIR) $(SRC_DIR)/*.proto

api: generateGrpc generateGateway generateSwagger