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

ALTER TABLE
    candle_one_minute
ADD
    CONSTRAINT candle_one_minute_pk PRIMARY KEY (ticker, "start");