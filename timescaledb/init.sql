CREATE TABLE IF NOT EXISTS temperature (
    time        TIMESTAMPTZ       NOT NULL,
    location    TEXT              NOT NULL,
    value       DOUBLE PRECISION  NULL
);

SELECT create_hypertable('temperature', 'time', if_not_exists => TRUE);
