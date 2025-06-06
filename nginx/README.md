# NGINX
A single proxy server that handles all the incoming connections, and forwards it to backend servers based on different paths

## Load Balancing
The upstream section defines list of servers for which the traffic will be routed in round robin fashion

## Workers
Usually there are as many workers as the number of CPUs. Each worker can handle `events.worker_connections` connections at a time. Workers run asynchronously to the master.

## Features
* Path rewrites
* Rate Limiter
* Doesn't support dynamic discovery of nodes. Have to integrate with Consul or Kubernetes to reload the config
