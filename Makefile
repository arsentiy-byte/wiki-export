run:
	go run ./cmd/app/main.go --config=./config/dev.yaml

mod:
	go mod tidy && go mod vendor
