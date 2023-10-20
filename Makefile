build:
	go build -o cmd/company/company cmd/company/main.go
lint:
	golangci-lint run --timeout=10m --disable-all -E misspell -E govet -E revive -E gofumpt -E gosec -E unparam -E goconst -E prealloc -E stylecheck -E unconvert -E errcheck -E ineffassign -E unused -E tparallel -E whitespace -E staticcheck -E gosimple -E gocritic
run:
	go run cmd/company/main.go --config=config/config.yml
migrate:
	go run cmd/migrate/main.go --config=config/config.yml
unit:
	go test ./... -race -count=1 -failfast
coverage:
	go test ./... -race -count=1 -failfast -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out