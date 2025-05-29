# Trino (Incomplete)
* Username admin/root at localhost:8080
* Executes distributed SQL queries across various data sources.

## Trino CLI
* Connect to MySQL `trino --server localhost:8080 --catalog mysql --schema demo`
* Connect to Postgres `trino --server localhost:8080 --catalog postgres --schema public`
* You can now create table and insert data as needed
* `SHOW CATALOGS;` to show the catalogs available
* `SHOW SCHEMAS FROM mysql;` to show the schemas
* `SHOW TABLES FROM mysql.demo;` to show tables
* Union two tables
```sql
select * from postgres.public.users UNION select * from  mysql.demo.users;
```

