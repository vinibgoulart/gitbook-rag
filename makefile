WORKER = cmd/worker/worker.go

worker:
	go run $(WORKER)

ai_cli:
	go run cmd/cli/cli.go

api:
	go run cmd/api/api.go