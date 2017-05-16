-- Data summary
SELECT * FROM
(SELECT count(*) StockCount from Stocks),
(SELECT count(*) PriceCount, MAX(SampleOn) LastPrice from Prices),
(SELECT count(*) AuthorCount, MAX(UpdatedOn) LastAuthorUpUpdate from Authors),
(SELECT count(*) PostCount, MAX(PostedOn) LastPost from Posts)
