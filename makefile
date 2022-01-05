name=app-center
version=v0.0.1
repository=lowcode
dockerHost=dockerhub.qingcloud.com
env:
#-- open go mod vendor --
	go mod vendor
test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

dev: linux
	./dev_auto.sh
linux:
	GOOS=linux GOOSARCH=amd64 go build -o $(name) ./cmd/.	
docker-test: env
	docker build -t $(dockerHost)/$(repository)/$(name):$(version) .
	docker push $(dockerHost)/$(repository)/$(name):$(version)
