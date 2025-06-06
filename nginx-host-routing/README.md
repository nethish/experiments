# NGINX Routing
Learnt that `localtest.me` simply resolves to `127.0.0.1` in our machine.
So, `api.localtest.me`, `notapi.localtest.me` is actually possible and different subdomains can actually be tested in local 

The `Host` header will be set to the domain name `Host: api.localtest.me`, and nginx can route the requests based on this header
