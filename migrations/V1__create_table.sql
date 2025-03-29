CREATE TABLE btc_usd_one_min (
  start_unix_time BIGINT PRIMARY KEY, --primary key indexes are automatically B+tree
  open_cents INT NOT NULL,
  high_cents INT NOT NULL,
  low_cents INT NOT NULL,
  close_cents INT NOT NULL,
  volume REAL NOT NULL,
  created_on TIMESTAMPTZ NOT NULL,
  
);

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_btc_usd_one_min_indicators (

);