DB_USER=admin
DB_PASS=admin
DB_HOST=localhost
DB_PORT=5432
DB_NAME=dumper
RECON_TICKER=500

PORT=8080

OUT_FILE=/Users/adrianbrad/workspace/go/dumper/dumper.out


run-test:
	OUT=./test.out DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) DB_USER=$(DB_USER) DB_PASS=$(DB_PASS) DB_NAME=$(DB_NAME) go test ./test -race -v

make run-with-closing-db:
	TIC=$(RECON_TICKER) PORT=$(PORT) OUT=$(OUT_FILE) DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) DB_USER=$(DB_USER) DB_PASS=$(DB_PASS) DB_NAME=$(DB_NAME) go run ./cmd/dumper/main.go -race