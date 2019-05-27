run:
	go run cmd/server/main.go

watch:
	ag -l | entr -r make run

tidy:
	go mod tidy
