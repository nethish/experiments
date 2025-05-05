$TTL 3600
@       IN  SOA ns1.example.com. admin.example.com. (
                2025050401 ; serial
                3600       ; refresh
                1800       ; retry
                604800     ; expire
                86400 )    ; minimum

@       IN  NS      ns1.example.com.
ns1     IN  A       127.0.0.1
@       IN  A       127.0.0.1
www     IN  A       127.0.0.1
