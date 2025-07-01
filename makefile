run:
	go run cmd/main.go

test:
	echo "running all of the tests"
	go test -coverprofile=coverage.out ./... 
	go tool cover -html=coverage.out -o coverage.html
