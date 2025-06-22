import socket
import ssl
import time

HOST = '127.0.0.1'
PORT = 8443

context = ssl.create_default_context()
context.check_hostname = False
context.verify_mode = ssl.CERT_NONE  # Don't verify for self-signed certs

# Create raw TCP socket
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Start timing
start = time.perf_counter()

# TCP connect
sock.connect((HOST, PORT))

# Wrap in TLS (TLS handshake happens here)
tls_sock = context.wrap_socket(sock, server_hostname=HOST)

# End timing
end = time.perf_counter()

# Total time = TCP + TLS handshake
print(f"Connection established in {(end - start) * 1000:.3f} ms")

# Optionally send/receive something
tls_sock.sendall(b"Hello")
response = tls_sock.recv(1024)
print(f"Server response: {response.decode(errors='ignore')}")

tls_sock.close()

