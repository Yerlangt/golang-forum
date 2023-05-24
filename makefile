check:
	go run ./cmd

build:
	docker build -t forum .
run:
	docker run -dp 8081:8081 --rm --name forum_container forum
stop:
	docker stop forum_container