# PostgreSQL Experiments
## Postgres Page Data
https://ketansingh.me/posts/how-postgres-stores-rows/

## Isolation Level
When there is an active transaction updating a data, postgres default isolation level will read data from previous version.

In MSSQL you can achive this by `DATABASE SET READ_COMMITTED_SNAPSHOT ON`

The default transaction isolation in postgres is Read Committed


