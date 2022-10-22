clean:
	rm -rf pb
	rm -rf swagger
	rm -rf tmp

gen:
	protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import ./proto/*.proto 

server1:
	go run cmd/server/main.go -port 50051

server2:
	go run cmd/server/main.go -port 50052

server1-tls:
	go run cmd/server/main.go -port 50051 -tls

server2-tls:
	go run cmd/server/main.go -port 50052 -tls

server:
	go run cmd/server/main.go -port 8080

server-tls:
	go run cmd/server/main.go -port 8080 -tls

rest:
	go run cmd/server/main.go -port 8081 -type rest -endpoint 0.0.0.0:8080

client:
	go run cmd/client/main.go -address 0.0.0.0:8080

client-tls:
	go run cmd/client/main.go -address 0.0.0.0:8080 -tls

test:
	rm -rf tmp
	mkdir tmp
	go test -cover -race -v ./...
	rm -rf tmp/*

cert:
	cd cert; ./gen.sh; cd ..

.PHONY: clean gen server client test cert 