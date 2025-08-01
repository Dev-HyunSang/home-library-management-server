dev:
	go run ./cmd/gocore/main.go dev
docker-run:
	docker compose --env-file ./config/.env.dev up
gen:
	go generate ./lib/ent