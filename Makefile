dev:
	go run ./cmd/gocore/main.go dev
docker-run:
	docker compose --env-file ./.env.dev up -d
gen:
	go generate ./lib/ent