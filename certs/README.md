
# create ca

```bash
openssl req -x509 -newkey rsa:4096 -days 3650 -nodes \
  -keyout ca.key.pem -out ca.cert.pem \
  -subj "/CN=My Custom CA"

```

# Create server key and csr

```bash
openssl req -newkey rsa:4096 -nodes -keyout server.key.pem -out server.csr.pem \
  -subj "/CN=localhost"
```

# Create server certificate

```bash
openssl x509 -req -in server.csr.pem -CA ca.cert.pem -CAkey ca.key.pem -CAcreateserial \
  -out server.cert.pem -days 365 \
  -extfile <(printf "subjectAltName=DNS:localhost,IP:127.0.0.1")
```

# Create client key and csr

```bash

openssl req -newkey rsa:2048 -nodes -keyout client.key.pem -out client.csr.pem \
  -subj "/CN=grpc-client"

openssl x509 -req -in client.csr.pem -CA ca.cert.pem -CAkey ca.key.pem -CAcreateserial \
  -out client.cert.pem -days 365

```
