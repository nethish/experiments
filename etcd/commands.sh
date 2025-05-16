docker exec -it etcd1 etcdctl member list
docker exec -it etcd1 etcdctl del "" --prefix
docker exec -it etcd1 etcdctl --endpoints=http://etcd1:2379,http://etcd2:2379,http://etcd3:2379 endpoint status --write-out=table
