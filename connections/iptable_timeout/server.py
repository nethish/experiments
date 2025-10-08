import socket
import time
from datetime import datetime

HOST = "0.0.0.0"
PORT = 5000

def log(msg):
    print(f"[{datetime.now().strftime('%H:%M:%S')}] {msg}", flush=True)

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.bind((HOST, PORT))
    s.listen(1)
    log(f"Server listening on {PORT}")
    conn, addr = s.accept()
    with conn:
        log(f"Connected by {addr}")
        while True:
            data = conn.recv(1024)
            if not data:
                log(f"Client disconnected")
                break
            message = data.decode().strip()
            log(f"Received: {message}")
            time.sleep(5)
            # Send a response back
            conn.sendall(f"ack: {message}\n".encode())

