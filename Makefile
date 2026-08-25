test:
	GOTOOLCHAIN=local go test ./... -count=1
race:
	GOTOOLCHAIN=local go test -race ./... -count=1
run:
	GOTOOLCHAIN=local go run ./cmd/server
