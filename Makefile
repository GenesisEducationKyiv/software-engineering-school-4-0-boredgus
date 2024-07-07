ENV_FILE=.env
DEV_ENV_FILE=.env.dev

start:
	ENV_FILE=${ENV_FILE} docker compose -f docker-compose.yaml --env-file ${ENV_FILE} up

start-dev:
	ENV_FILE=${DEV_ENV_FILE} docker compose -f docker-compose-dev.yaml  --env-file ${DEV_ENV_FILE} up

generate-mocks:
	cd gateway && make generate-mocks;
	cd service/currency && make generate-mocks;
	cd service/dispatch && make generate-mocks;
	cd service/notification && make generate-mocks;

lint:
	golangci-lint run -c ./gateway/.golangci.yaml ./gateway/...
	golangci-lint run -c ./service/dispatch/.golangci.yaml ./service/dispatch/...
	golangci-lint run -c ./service/currency/.golangci.yaml ./service/currency/...
	golangci-lint run -c ./service/notification/.golangci.yaml ./service/notification/... 
		

test:
	go test \
		./service/dispatch/... \
		./service/currency/... \
		./service/notification/... \
		./gateway/... \
		-coverprofile="test-coverage.txt" \
		-covermode count

	go tool cover -func="test-coverage.txt"

test-coverage:
	go tool cover -html="test-coverage.txt"


currency_service_proto=contracts/proto/services/currency_service.proto
dispatch_service_proto=contracts/proto/services/dispatch_service.proto

messages_proto=contracts/proto/messages/messages.proto
dispatch_messages_proto=contracts/proto/messages/dispatch/dispatch_messages.proto ${messages_proto}
subscription_messages_proto=contracts/proto/messages/subscription/subscription_messages.proto ${messages_proto}

generate-proto:
# for gateway
	protoc --go_out=./gateway/internal/clients/currency/gen  \
		--go-grpc_out=./gateway/internal/clients/currency/gen \
		--proto_path=contracts/proto \
		${currency_service_proto}

	protoc --go_out=./gateway/internal/clients/dispatch/gen  \
		--go-grpc_out=./gateway/internal/clients/dispatch/gen \
		--proto_path=contracts/proto \
		${dispatch_service_proto}

# for currency service
	protoc --go_out=./service/currency/internal/grpc/gen \
		--go-grpc_out=./service/currency/internal/grpc/gen \
		--proto_path=contracts/proto \
		${currency_service_proto}

# for dispatch service
	protoc --go_out=./service/dispatch/internal/grpc/gen \
		--go-grpc_out=./service/dispatch/internal/grpc/gen \
		--proto_path=contracts/proto \
		${dispatch_service_proto}

	protoc --go_out=./service/dispatch/internal/broker/gen \
		--proto_path=contracts/proto \
		${subscription_messages_proto}

# for notification service
	protoc --go_out=./service/notification/internal/broker/gen \
		--proto_path=contracts/proto \
		${subscription_messages_proto}

	protoc --go_out=./service/notification/internal/broker/gen \
		--proto_path=contracts/proto \
		${dispatch_messages_proto}

	protoc --go_out=./service/notification/internal/clients/currency/gen  \
		--go-grpc_out=./service/notification/internal/clients/currency/gen \
		--proto_path=contracts/proto \
		${currency_service_proto}
