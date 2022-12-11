deps:
	go mod download
	go mod tidy

run:
	go run .

dotenv:
	export $(grep -v '^#' .env | xargs)