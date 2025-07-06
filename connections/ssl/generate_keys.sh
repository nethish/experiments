# Generate a private key
openssl genrsa -out server.key 2048

# Generate a self-signed certificate
openssl req -new -x509 -key server.key -out server.crt -days 3650 -subj "/CN=localhost"

# Combine key and cert into a .pem file (often used by servers)
cat server.key server.crt >server.pem

# Secure permissions for the key
chmod 600 server.key server.pem
