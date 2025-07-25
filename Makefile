up:
	docker compose down
	docker compose up --build -d
	docker image prune -f

ps:
	docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}"

test:
	go test ./internal/application/service -coverprofile=coverage.out 

cover:
	go tool cover -html=coverage.out

log:
	git log --oneline -n 10