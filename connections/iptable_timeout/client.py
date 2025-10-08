import socket
import time
from datetime import datetime

SERVER_HOST = "tcp_server"
SERVER_PORT = 5000

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.settimeout(30)  # read timeout
sock.connect((SERVER_HOST, SERVER_PORT))

def log(msg):
    print(f"[{datetime.now().strftime('%H:%M:%S')}] {msg}", flush=True)

log("Connected to server")


try:
    counter = 1
    while True:
        msg = f"message {counter}"
        sock.sendall(msg.encode())
        log(f"Sent: {msg}")

        # Wait for server response (recv will block and respect timeout)
        data = sock.recv(1024)
        log(f"Received: {data.decode().strip()}")

        time.sleep(1)
        counter += 1
except socket.timeout:
    log("Timeout reached! Server not responding.")
except Exception as e:
    log(f"Exception: {e}")
finally:
    sock.close()
    log("Connection closed")

