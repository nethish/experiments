import socket

HOST = "127.0.0.1"
PORT = 8080
NUM_CONNECTIONS = 10000  # try more than backlog

sockets = []

for i in range(NUM_CONNECTIONS):
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.settimeout(1)  # timeout to avoid hanging forever
        s.connect((HOST, PORT))
        sockets.append(s)
        print(f"Connection {i + 1} succeeded")
    except Exception as e:
        print(f"Connection {i + 1} failed with error: {e}")
        break

print(f"Total successful connections: {len(sockets)}")

# Keep connections open to hold the queue
input("Press Enter to exit and close all connections...")

# Cleanup
for s in sockets:
    s.close()
