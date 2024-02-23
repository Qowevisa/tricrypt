def: server client
	@

all: rm def
	@

rm:
	rm ./bin/* 2>/dev/null || true

server: server.crt server.key
	go build -o ./bin/$@ ./cmd/$@

client: ca.crt
	go build -o ./bin/$@ ./cmd/$@

gen_test_certs:
	openssl ecparam -genkey -name prime256v1 -out server.key
	openssl req -new -x509 -key server.key -out server.pem -days 3650

gen_certs:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -config san.cnf
	#openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365

all: ca.crt server.crt client.crt

ca.key:
	openssl genrsa -out ca.key 4096

ca.crt: ca.key
	openssl req -new -x509 -days 365 -key ca.key -out ca.crt -subj "/C=US/ST=YourState/L=YourCity/O=YourOrganization/CN=YourCA"

server.key:
	openssl genrsa -out server.key 4096

server.csr: server.key
	openssl req -new -key server.key -out server.csr -subj "/C=US/ST=YourState/L=YourCity/O=YourOrganization/CN=server.yourdomain.com"

server.crt: server.csr ca.crt ca.key
	openssl x509 -req -days 365 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt

client.key:
	openssl genrsa -out client.key 4096

client.csr: client.key
	openssl req -new -key client.key -out client.csr -subj "/C=US/ST=YourState/L=YourCity/O=YourOrganization/CN=client.yourdomain.com"

client.crt: client.csr ca.crt ca.key
	openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt

clean:
	rm -f ca.key ca.crt server.key server.csr server.crt client.key client.csr client.crt

