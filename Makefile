SRC_DIR=server/api
DST_DIR=$(SRC_DIR)
GATEWAY_VERSION=v1.12.2-0.20191203171358-3c06998610d4
generateGrpc:
	protoc -I=$(SRC_DIR) \
		-I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$(GATEWAY_VERSION)/third_party/googleapis \
    	-I=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$(GATEWAY_VERSION) \
		--go_out=plugins=grpc:$(DST_DIR) \
		$(SRC_DIR)/*.proto
generateGateway:
	protoc -I=$(SRC_DIR) \
		-I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$(GATEWAY_VERSION)/third_party/googleapis \
    	-I=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$(GATEWAY_VERSION) \
		--grpc-gateway_out=logtostderr=true:$(DST_DIR) \
		$(SRC_DIR)/*.proto

generateSwagger:
	protoc -I=$(SRC_DIR) \
    		-I=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$(GATEWAY_VERSION)/third_party/googleapis \
    		-I=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$(GATEWAY_VERSION) \
    		--swagger_out=logtostderr=true,allow_merge=true:$(DST_DIR) $(SRC_DIR)/*.proto

api: generateGrpc generateGateway generateSwagger