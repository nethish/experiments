# Stress test MongoDB
* Write 50 million data points to a collection
* Don't have an index
* Do high concurrent reads on the collection
* View the number of read and write tickets available with `db.serverStatus().wiredTiger.concurrentTransactions`
  * It gets stuck at 28 for some reason
* Only MongoDB 6 versions have default read and write tickets set to 128. Later version have a good algo to dynamically adjust this

## Run
* Uncomment `seedData` method only once
```
docker run -d \
  --name mongo6 \
  -p 27017:27017 \
  mongo:6
```

Example output
```
test_db> db.serverStatus().wiredTiger.concurrentTransactions
{
  write: { out: 1, available: 127, totalTickets: 128 },
  read: { out: 7, available: 121, totalTickets: 128 }
}
```
