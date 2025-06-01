# Experiment with connections
* Theoretically no limit on number of connections. A connection `[serverIP:serverPort:clientIP:clientPort]`
* The Python example is more simpler than the go example

```bash
# Increase the size of the connection queue - Apple Silicon
sudo sysctl -w kern.ipc.somaxconn=10000
```
