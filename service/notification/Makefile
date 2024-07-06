
start:
	go run cmd/server/main.go

test: 
	go test ./... -coverprofile="test-coverage.txt" -covermode count
	go tool cover -func="test-coverage.txt"

generate-mocks:
	mockery --config=./.mockery.yaml
	