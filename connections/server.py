import socket
import time

HOST = "127.0.0.1"
PORT = 8080
BACKLOG = 9000

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server:
    server.bind((HOST, PORT))
    server.listen(BACKLOG)
    print(f"Server listening on {HOST}:{PORT} with backlog={BACKLOG}")

    # Do NOT accept connections
    print("Not accepting connections, just sleeping to let queue fill up...")
    while True:
        time.sleep(10)
