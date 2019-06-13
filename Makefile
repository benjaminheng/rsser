dockerbuild:
	docker build . -t rsser

dockerrun:
	docker rm rsser || true
	docker run --rm --name rsser -p 8100:8000 rsser

run:
	go run cmd/server/main.go

watch:
	ag -l | entr -r make run

tidy:
	go mod tidy
