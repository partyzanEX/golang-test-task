deps:
	glide install
gotask:
	CGO_ENABLED=0 GOOS=linux go build -o gotask ./server/main.go
build:
	docker build --no-cache -t partyzanex/golang-test .
run:
	docker run --name gotask -p "3030:3030" -d partyzanex/golang-test
clean:
	rm ./gotask
