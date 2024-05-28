run:
	go run .

run-docker:
	docker build -t basic-forms .
	docker run -p 8080:8080 --env-file .env --env GIN_MODE=release basic-forms

deploy:
	caprover deploy --default