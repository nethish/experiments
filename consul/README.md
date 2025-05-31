# Consul
This repository demonstrates service discovery example
Consul runs at `http://localhost:8500/`
Hit any app with curl `"http://localhost:8080/call?service=app3"`

## Usecases
* KV Store
* First class service discovery
* Health Checks and monitoring

## How it differs from Zookeeper
* First class leader election
* Stores config via znodes
* Distributed Locks

## Diff
* Both can do what the other can do but for few things one is better than the other

