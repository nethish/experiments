# tls_server.py
# openssl req -x509 -newkey rsa:2048 -nodes -keyout server.key -out server.crt -days 365
# openssl s_client -connect localhost:8443


import socket
import ssl

HOST = "127.0.0.1"
PORT = 8443

context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
context.load_cert_chain(certfile="server.crt", keyfile="server.key")

with socket.socket(socket.AF_INET, socket.SOCK_STREAM, 0) as sock:
    sock.bind((HOST, PORT))
    sock.listen(5)
    print(f"Listening on {HOST}:{PORT} with TLS...")

    with context.wrap_socket(sock, server_side=True) as ssock:
        conn, addr = ssock.accept()
        print(f"Connection from {addr}")
        data = conn.recv(1024)
        print(f"Received: {data.decode(errors='ignore')}")
        conn.sendall(b"Hello over TLS")
        conn.close()
