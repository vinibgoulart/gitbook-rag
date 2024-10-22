SERVER_PATH = cmd/server/server.go

server:
	go run $(SERVER_PATH)

ai_cli:
	go run cmd/cli/cli.go

api:
	go run cmd/api/api.go