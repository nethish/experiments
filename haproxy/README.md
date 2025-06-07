# HAProxy
High Availability Proxy

## Features
* Health checks
* Load balancing - least conn, source ip, uri, random
* TLS Termination
* ACL
* Effective connection reuse for backend servers
* Custom Lua scripts

## As HTTP Proxy
It can act as either L4 or L7 proxy. For L7 proxy use `haproxy.cfg`. This will only forward http requests to the server

When you configure as http proxy, you can see the stats at 
```
curl -v http://localhost:8080
http://localhost:8404/stats
```

## As TCP Proxy
Use the `l4.cfg` to make it a L4 proxy. It will blindly forward traffic to the backend webservers

```bash
nc localhost 8080
GET / HTTP/1.1
Host: localhost
```

Note: Nginx can also act as L4 proxy when you use `stream` context, but that's like addon.
