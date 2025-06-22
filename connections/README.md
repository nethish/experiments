# Experiment with connections
* Theoretically no limit on number of connections. A connection is `[serverIP:serverPort:clientIP:clientPort]`
* The Python example is more simpler than the go example
* The root folder python example demonstrates what happens if number of accept is exhausted (backlog queue size). The connection will timeout if server doesn't accept the connection.

## SSL Folder
* The root layer server-client example demonstrates how quickly a server can accept a tcp connection with TLS. Locally it doesn't even take 1 millisecond
* Meanwhile, with TLS, the connection could take 5 millisecond to get and verify certs

```bash
python client.py
# Connection established in 5.679 ms
```


## System level settings for connection tuning
* Maximum backlog queue size
```bash
# Increase the size of the connection queue - Apple Silicon
sudo sysctl -w kern.ipc.somaxconn=10000

# Prints the value - Default 128
sysctl kern.ipc.somaxconn
```
