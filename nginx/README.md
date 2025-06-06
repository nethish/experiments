# NGINX
A single proxy server that handles all the incoming connections, and forwards it to backend servers based on different paths

Nginx usually streams the response to the downstream servers. The response end is either marked using `Content-Length: 100` header or `Transfer Encoding: chucked` marks it's end with empty chunk

## Load Balancing
The upstream section defines list of servers for which the traffic will be routed in round robin fashion

## Workers
Usually there are as many workers as the number of CPUs. Each worker can handle `events.worker_connections` connections at a time. Workers run asynchronously to the master.

## Multiple Servers
Learnt that `localtest.me` simply resolves to `127.0.0.1` in our machine.
So, `api.localtest.me`, `notapi.localtest.me` is actually possible and different subdomains can actually be tested in local 

The `Host` header will be set to the domain name `Host: api.localtest.me`, and nginx can route the requests based on this header. The Host header based routing is part of HTTP1.1 spec

Quote the HTTP 1.1 RFC

```
The "Host" header field in a request provides the host and port information from the target URI, enabling the origin server to distinguish among resources while servicing requests for multiple host names.
```

The below URLs will all resolve correctly

```bash
http://api.localtest.me/ # This only goes to app_1
http://frontend.localtest.me/
http://localhost/       # Goes to frontend
http://localhost/api/   # This round robins app_1 and app_2
```


The below nginx config is the catch all block
```
server_name _;  # catch-all or fallback

```

## Features
* Path rewrites
* Rate Limiter
* Doesn't support dynamic discovery of nodes. Have to integrate with Consul or Kubernetes to reload the config
