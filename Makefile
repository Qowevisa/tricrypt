def: server client
	@

all: rm def
	@

rm:
	rm ./bin/* 2>/dev/null || true

server:
	go build -o ./bin/$@ ./cmd/$@

client:
	go build -o ./bin/$@ ./cmd/$@

gen_certs:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -config san.cnf
	#openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365

