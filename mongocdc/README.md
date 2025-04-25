# MongoDB CDC Stream 

## About Mongo CDC
* CDC Stream is a feature of Mongo ReplicaSet
* The MongoDB should run in replica mode (even if it's single instance) for the CDC to work.
* The MongoDB writes all operations to oplog, and streams it from there.
* You can resume the stream from last consumed event using `resume_after` field. See `main.py` for example
* If the oplog expires due to space exhaustion or high write rate, then resume after don't work
* Below is an example from local
```bash
myReplicaSet [direct: primary] test> rs.printReplicationInfo()
actual oplog size
'1263.517822265625 MB'
---
configured oplog size
'1263.517822265625 MB'
---
log length start to end
'15436 secs (4.29 hrs)'
---
oplog first event time
'Tue Apr 15 2025 16:26:14 GMT+0000 (Coordinated Universal Time)'
---
oplog last event time
'Tue Apr 15 2025 20:43:30 GMT+0000 (Coordinated Universal Time)'
---
now
'Tue Apr 15 2025 20:43:35 GMT+0000 (Coordinated Universal Time)'

```

## Python
* Install pymongo
* Start up the server
* Insert docs 
```
docker exec -it mongo1 mongosh

use test

db.users.insertOne({name: "nethi", age: 25})
```

* The change events should be reflected in the python output

## Go
* `go mod tidy` for deps
* `go run main.go`
* Same as above

## References
Thanks to the [medium post](https://alek-cora-glez.medium.com/change-data-capture-service-for-mongodb-a87634d2b16b)

