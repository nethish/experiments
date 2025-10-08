# Connection timeout simulation

Run the `run.sh` file to spin up client and server. The client sends a message to the server, and the server waits for 5 seconds before replying.

We block the client by disconnecting it from the network. Then the client won't receive any packets, but times out after 30 seconds.

```bash
# This was run in mac
docker network disconnect iptable_timeout_default tcp_client
docker network connect iptable_timeout_default tcp_client
```
