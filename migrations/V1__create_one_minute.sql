-- one minute candles that will be used to create other timeframes/indicators 
CREATE TABLE candle_one_minute (
    ticker VARCHAR(16) NOT NULL,
    "start" BIGINT NOT NULL,
    "open" DECIMAL(18, 8) NOT NULL,
    high DECIMAL(18, 8) NOT NULL,
    low DECIMAL(18, 8) NOT NULL,
    "close" DECIMAL(18, 8) NOT NULL,
    volume DOUBLE PRECISION NOT NULL
);

ALTER TABLE candle_one_minute
ADD CONSTRAINT candle_one_minute_pk PRIMARY KEY (ticker, "start");

-- create hypertable with 6 month chunks
SELECT create_hypertable('candle_one_minute', by_range('start', 15778440));

/*
-- check if hypertable has been made
SELECT *
FROM timescaledb_information.hypertables
*/

/*
-- timescale apache license doesn't support columnar compression
ALTER TABLE candle_one_minute SET (
    timescaledb.compress,
    timescaledb.compress_orderby = 'start ASC',
    timescaledb.compress_segmentby = 'ticker'
);
*/