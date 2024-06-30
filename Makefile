ENV_FILE=.env

start:
	ENV_FILE=${ENV_FILE} docker compose -f docker-compose.yaml --env-file ${ENV_FILE} up

generate-mocks:
	mockery --config=config/.mockery.yaml

lint:
	golangci-lint run -c .golangci.yaml \
		./service/dispatch/... \
		./service/currency/... \
		./daemon/dispatch/... \
		./gateway/...

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
