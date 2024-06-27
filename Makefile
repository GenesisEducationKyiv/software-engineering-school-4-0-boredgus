ENV_FILE=.env

start:
	ENV_FILE=${ENV_FILE} docker compose -f docker-compose.yaml --env-file ${ENV_FILE} up

generate-mocks:
	mockery --config=config/.mockery.yaml

lint:
	golangci-lint run -c .golangci.yaml

test:
	go test ./... -coverprofile="test-coverage.txt" -covermode count
	go tool cover -func="test-coverage.txt"

test-coverage:
	go tool cover -html="test-coverage.txt"

generate-grpc:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/services/currency/grpc/currencyService.proto \
		internal/services/dispatch/grpc/dispatchService.proto
