#!/bin/bash

# Connect to HiveServer2 and run demo queries
docker exec -it hiveserver2 beeline -u jdbc:hive2://localhost:10000 -e "

-- Create a sample database
CREATE DATABASE IF NOT EXISTS demo;
USE demo;

-- Create a sample table
CREATE TABLE IF NOT EXISTS employees (
    id INT,
    name STRING,
    department STRING,
    salary DOUBLE,
    hire_date DATE
)
STORED AS TEXTFILE;

-- Insert sample data
INSERT INTO employees VALUES
(1, 'John Doe', 'Engineering', 75000.0, '2020-01-15'),
(2, 'Jane Smith', 'Marketing', 65000.0, '2019-03-22'),
(3, 'Bob Johnson', 'Engineering', 80000.0, '2021-06-10'),
(4, 'Alice Brown', 'HR', 55000.0, '2018-11-05'),
(5, 'Charlie Davis', 'Engineering', 90000.0, '2022-02-28');

-- Query examples
SELECT 'Total employees:' as metric, COUNT(*) as value FROM employees
UNION ALL
SELECT 'Average salary:' as metric, ROUND(AVG(salary), 2) as value FROM employees;

-- Department analysis
SELECT 
    department,
    COUNT(*) as employee_count,
    ROUND(AVG(salary), 2) as avg_salary,
    MAX(salary) as max_salary
FROM employees 
GROUP BY department 
ORDER BY avg_salary DESC;

-- Show tables and describe structure
SHOW TABLES;
DESCRIBE employees;
"
