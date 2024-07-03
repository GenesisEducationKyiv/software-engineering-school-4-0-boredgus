ENV_FILE=.env

start:
	ENV_FILE=${ENV_FILE} docker compose -f docker-compose.yaml --env-file ${ENV_FILE} up

generate-mocks:
	mockery --config=./gateway/.mockery.yaml
	mockery --config=./service/currency/.mockery.yaml
	mockery --config=./service/dispatch/.mockery.yaml

lint:
	golangci-lint run -c ./gateway/.golangci.yaml ./gateway/...
	golangci-lint run -c ./service/dispatch/.golangci.yaml ./service/dispatch/...
	golangci-lint run -c ./service/currency/.golangci.yaml ./service/currency/...
	golangci-lint run -c ./daemon/dispatch/.golangci.yaml ./daemon/dispatch/... 
		

test:
	go test \
		./service/dispatch/... \
		./service/currency/... \
		./daemon/dispatch/... \
		./gateway/... \
		-coverprofile="test-coverage.txt" \
		-covermode count

	go tool cover -func="test-coverage.txt"

test-coverage:
	go tool cover -html="test-coverage.txt"

generate-grpc:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		service/currency/internal/grpc/gen/currency_service.proto \
		service/dispatch/internal/grpc/gen/dispatch_service.proto

	protoc --go_out=./service/dispatch/internal/clients/currency/gen  \
		--go-grpc_out=./service/dispatch/internal/clients/currency/gen \
		service/currency/internal/grpc/gen/currency_service.proto

	protoc --go_out=./daemon/dispatch/internal/clients/dispatch/gen  \
		--go-grpc_out=./daemon/dispatch/internal/clients/dispatch/gen \
		service/dispatch/internal/grpc/gen/dispatch_service.proto
	
	protoc --go_out=./gateway/internal/clients/currency/gen  \
		--go-grpc_out=./gateway/internal/clients/currency/gen \
		service/currency/internal/grpc/gen/currency_service.proto
	protoc --go_out=./gateway/internal/clients/dispatch/gen  \
		--go-grpc_out=./gateway/internal/clients/dispatch/gen \
		service/dispatch/internal/grpc/gen/dispatch_service.proto
