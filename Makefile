run:
	go run ./...

alive:
	curl -s http://127.0.0.1:9090/api/v1/

post:
	curl -X POST -H "Content-Type: application/json" -d '{"name": "linuxize", "agency": "3", "account": "5"}' http://localhost:9090/api/v1/bank

get:
	curl -s http://localhost:9090/api/v1/bank
