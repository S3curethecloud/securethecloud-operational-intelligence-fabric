.PHONY: up down reset api ai frontend event test opa-test

up:
	docker compose up -d

down:
	docker compose down

reset:
	docker compose down -v
	docker compose up -d

api:
	cd api && go run ./cmd/server

ai:
	cd ai-service && python3 -m venv .venv && . .venv/bin/activate && pip install -e . && uvicorn app.main:app --host 0.0.0.0 --port 8081 --reload

frontend:
	cd frontend && npm install && npm run dev

event:
	./scripts/simulate_runtime_event.sh

opa-test:
	docker run --rm -v "$$(pwd)/policy/opa:/policies" openpolicyagent/opa:latest test /policies

test:
	cd api && go test ./...
	cd ai-service && python3 -m pytest tests || true
