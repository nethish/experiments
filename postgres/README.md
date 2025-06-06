# PostgreSQL Experiments
## Postgres Page Data
https://ketansingh.me/posts/how-postgres-stores-rows/

## Isolation Level
When there is an active transaction updating a data, postgres default isolation level will read data from previous version.
This is called MVCC. When a row is updated, it's MARK DELETE + INSERT with monotonically increasing txn id.
When a transaction with id txn_id is executed, it can only see rows whos max_tid (part of row tuple) is less than txn_id.

In MSSQL you can achive this by `DATABASE SET READ_COMMITTED_SNAPSHOT ON`

The default transaction isolation in postgres is Read Committed


