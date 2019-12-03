help:
	@echo "You can perform the following:"
	@echo ""
	@echo "  check         Format, lint, vet, and test Go code"
	@echo "  cover         Show test coverage in html"
	@echo "  deploy        Deploy to IBM Cloud Foundry"
	@echo "  dev           Build and run for local development OS"
	@echo "  local         Build for local development OS"
	@echo "  pg            Start postgres in Docker on port 5432"
	@echo "  psql          Connect using psql (password: docker)"
	@echo "  delete        Delete database"

check:
	@echo 'Formatting, linting, vetting, and testing Go code'
	go fmt ./...
	golint ./...
	go vet ./...
	go test ./... -cover

cover:
	@echo 'Test coverage in html'
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

#  Compile the project to run locally on your machine
local:
	go build -o dist/lenslocked .

dev: local
	dist/lenslocked

deploy:
	go mod tidy
	gcloud builds submit --tag gcr.io/todobackendgcr/todobackend-gcr
	gcloud run deploy --image gcr.io/todobackendgcr/todobackend-gcr --platform managed

pg:
	docker run --rm --name pg-docker -e POSTGRES_PASSWORD=docker -d -p 5432:5432 -v ~/docker/volumes/postgres:/var/lib/postgresql/data postgres

psql:
	psql -h localhost -U postgres -d postgres

delete:
	go build -o dist/delete cmd/delete/*
	dist/delete
