-- delete old rows once a month
SELECT cron.schedule('0 3 1 * *', $$DELETE FROM candle_one_minute WHERE "start" < extract(epoch from now() - interval '3 year')$$);
