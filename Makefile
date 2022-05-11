PLATFORM=$(shell uname -m)
DATETIME=$(shell date "+%Y%m%d%H%M%S")
VERSION=v0.0.1-SNAPSHOT

release:
	@cd front && npm run build -- prod && cp -af dist/* ../service/src/resources/
	@cd service/src && GOPATH=${GOPATH} go build -o ../bin/server.bin

front-debugging:
	@cd front && npm run dev

backend-debugging:
	@cd service/src && GOPATH=${GOPATH} go build -gcflags=all="-N -l" -o ../bin/server.bin

clean:
	@rm -rf ./bin

docker:
	@docker build . -t harbor.timeforward.cn:8443/public/management-backend:$(VERSION)-$(DATETIME)

