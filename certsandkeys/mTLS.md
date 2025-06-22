# mTLS Demo
```bash
# Generate RootCA
openssl genrsa -out rootCA/rootCA.key.pem 4096

# Generate RootCA Certificate
openssl req -x509 -new -nodes -key rootCA/rootCA.key.pem \
  -sha256 -days 3650 -out rootCA/rootCA.cert.pem \
  -subj "/C=IN/ST=TN/L=Chennai/O=ExampleRootCA/CN=example-root"

# Intermediary CA
openssl genrsa -out intermediateCA/intermediate.key.pem 4096

# Generate intermeidate CA CSR
openssl req -new -key intermediateCA/intermediate.key.pem \
  -out intermediateCA/intermediate.csr.pem \
  -subj "/C=IN/ST=TN/L=Chennai/O=ExampleIntermediateCA/CN=example-intermediate"

# Sign intermediate CA with root ca
openssl x509 -req -in intermediateCA/intermediate.csr.pem \
  -CA rootCA/rootCA.cert.pem -CAkey rootCA/rootCA.key.pem \
  -CAcreateserial -out intermediateCA/intermediate.cert.pem \
  -days 1825 -sha256

#=========================

# Server certificate
openssl genrsa -out server/server.key.pem 2048

# Server cert CSR
openssl req -new -key server/server.key.pem -out server/server.csr.pem \
  -subj "/C=IN/ST=TN/L=Chennai/O=ExampleServer/CN=localhost"

# Server cert with intermediate CA
openssl x509 -req -in server/server.csr.pem \
  -CA intermediateCA/intermediate.cert.pem \
  -CAkey intermediateCA/intermediate.key.pem \
  -CAcreateserial -out server/server.cert.pem \
  -days 825 -sha256


#===============================

# Client certificate 
openssl genrsa -out client/client.key.pem 2048

# Client CSR
openssl req -new -key client/client.key.pem -out client/client.csr.pem \
  -subj "/C=IN/ST=TN/L=Chennai/O=ExampleClient/CN=client.local"

# Client cert
openssl x509 -req -in client/client.csr.pem \
  -CA intermediateCA/intermediate.cert.pem \
  -CAkey intermediateCA/intermediate.key.pem \
  -CAcreateserial -out client/client.cert.pem \
  -days 825 -sha256


# Create CA Bundle
cat intermediateCA/intermediate.cert.pem rootCA/rootCA.cert.pem > ca-chain.cert.pem

# Verify server cert
openssl verify -CAfile ca-chain.cert.pem server/server.cert.pem

# Verify client cert
openssl verify -CAfile ca-chain.cert.pem client/client.cert.pem

# =========================================

openssl s_server -accept 8443 \
  -cert server/server.cert.pem \
  -key server/server.key.pem \
  -CAfile ca-chain.cert.pem \
  -Verify 1

# Verify 1 - this will now verify the client certificate as well


openssl s_client -connect localhost:8443 \
  -cert client/client.cert.pem \
  -key client/client.key.pem \
  -CAfile ca-chain.cert.pem
```
