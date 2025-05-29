CREATE TABLE users (id int, name varchar);
INSERT INTO users VALUES (1, 'nethish');
INSERT INTO users VALUES (2, 'sash');

---


-- POSTGRES
CREATE TABLE customers (id int, name varchar, region varchar);
INSERT INTO customers VALUES (1, 'Alice', 'West');
INSERT INTO customers VALUES (2, 'Bob', 'East');
INSERT INTO customers VALUES (3, 'Carol', 'South');
INSERT INTO customers VALUES (4, 'Dave', 'North');
INSERT INTO customers VALUES (5, 'Nethish', 'Coimbatore');

-- MYSQL
CREATE TABLE sales (id int, customer_id int, amount int);
INSERT INTO sales VALUES (1, 1, 100);
INSERT INTO sales VALUES (2, 1, 200);
INSERT INTO sales VALUES (3, 2, 300);
INSERT INTO sales VALUES (4, 3, 400);
INSERT INTO sales VALUES (5, 5, 500);


SELECT c.region, SUM(s.amount) AS total_sales
FROM postgres.public.customers c
JOIN mysql.demo.sales s ON c.id = s.customer_id
GROUP BY c.region;

