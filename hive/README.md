# Hive




## Create a database
* Initialize metastore first
```bash
docker exec -it hive-metastore bash
schematool -initSchema -dbType postgres
```

* Exec into the hiveserver and connect using beeline or hive
```bash
docker exec -it hiveserver bash
beeline -u jdbc:hive2://hiveserver:10000/default -n hiveuser -p hivepassword
# You can optionally use `hive` command
```

* Database commands
```SQL
SHOW DATABASES;
CREATE DATABASE first_database;
USE first_database;

CREATE TABLE users (
  id INT,
  fruit STRING
)
ROW FORMAT DELIMITED
FIELDS TERMINATED BY '\t'
STORED AS TEXTFILE;

-- Don't forget to copy the kv1.txt available locally
LOAD DATA LOCAL INPATH '/opt/hive/data/kv1.txt' OVERWRITE INTO TABLE users;

SELECT * FROM users;


```
