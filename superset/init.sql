CREATE TABLE IF NOT EXISTS sales (
    id SERIAL PRIMARY KEY,
    product VARCHAR(50),
    region VARCHAR(50),
    sales_date DATE,
    quantity INTEGER,
    price NUMERIC(10,2)
);

INSERT INTO sales (product, region, sales_date, quantity, price) VALUES
('Product A', 'North', '2025-01-01', 10, 25.00),
('Product B', 'South', '2025-01-01', 15, 30.00),
('Product A', 'East', '2025-02-01', 20, 25.00),
('Product C', 'West', '2025-02-15', 5, 50.00),
('Product B', 'North', '2025-03-01', 8, 30.00),
('Product A', 'South', '2025-03-05', 12, 25.00);
