run: build
	@./bin/dreampicai

install:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download
	@npm install -D tailwindcss
	@npm install -D daisyui@latest

css:
	@npx tailwindcss -i view/css/app.css -o public/styles.css

templ:
	@templ generate --proxy=http://localhost:3000

build:
	@templ generate view
	@go build -tags dev -o bin/dreampicai main.go 


migration: ## Migrations against the database
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

seed:
	@go run cmd/seed/main.go