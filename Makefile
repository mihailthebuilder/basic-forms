run:
	go run .

run-reset:
	go run . -reset

run-docker:
	docker build -t basic-forms .
	docker run -p 8080:8080 --env-file .env --env GIN_MODE=release basic-forms