# Some materials on certificats and keys
* `X.509` is a format. Contains issues to and issued by, with public key
* PEM is the file encoding format.
  * It has the BEGIN and END rules. CERTIFICATE/ RSA PRIVATE KEY etc
* When generating a X.509 certificate, you get a cert and a private key
* Then extract the public key from the certificate by
```bash
openssl x509 -in your_certificate.cer -pubkey -noout > public_key.pem
```

## Sites
* Generate a `X.509` key using https://www.samltool.com/self_signed_certs.php
* Decode it https://certificatedecoder.dev/

## Some useful commands

```bash
# Generate RSA Key
openssl genrsa -out private_key.pem 2048

# The public key is embedded in it, you can run below to get the public key
openssl rsa -in private_key.pem -pubout -out public_key.pem

# Generate EC private key
openssl ecparam -name prime256v1 -genkey -noout -out ec_private_key.pem

# Public Key
openssl ec -in ec_private_key.pem -pubout -out ec_public_key.pem

# ----------------------------

# Self signed cert
openssl req -x509 -new -nodes -key private_key.pem -sha256 -days 365 -out cert.pem

# View the contents of the certificate
openssl x509 -in cert.pem -text -noout

# Contents of the private key
openssl rsa -in private_key.pem -check

# Convert PEM to DER
openssl x509 -in cert.pem -outform der -out cert.der

# Convert DER to PEM
openssl x509 -in cert.der -inform der -out cert.pem

# PEM to PKCS#12 (includes cert + key)
# Bundling is a technique used to reduce file juggling (private key, certificate etc) in the server app
# Windows, and Java Keystore (java's bundle format) widely supports it
# PKCS is bundle format and can store
#   * private key
#   * multiple certs
#   * CA chain
openssl pkcs12 -export -out cert_bundle.p12 -inkey private_key.pem -in cert.pem
```


## SSL TLS 
```bash
# Test SSL/ TLS on server
openssl s_client -connect google.com:443

# Show supported ciphers
openssl ciphers -v
```

## Popular PKCS
* PKCS#1 RSA Cryptography Standard
* PKCS#3 Diffie-Hellman Key Exchange
* PKCS#5 Password-Based Encryption (PBE)
* PKCS#7 Cryptographic Message Syntax (CMS)
* PKCS#8 Private Key Information Syntax
* PKCS#10 Certificate Signing Request (CSR)
* PKCS#11 Cryptographic Token Interface (API)
* PKCS#12 Personal Information Exchange Syntax
* PKCS#15 Cryptographic Token Information Format

