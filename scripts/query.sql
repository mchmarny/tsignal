-- Data summary
SELECT * FROM
(SELECT count(*) StockCount from Stocks),
(SELECT count(*) PriceCount, MAX(SampleOn) LastPrice from Prices),
(SELECT count(*) AuthorCount, MAX(UpdatedOn) LastAuthorUpUpdate from Authors),
(SELECT count(*) PostCount, MAX(PostedOn) LastPost from Posts)

-- Posts per Stock
SELECT Symbol, COUNT(*) Posts
FROM Posts
GROUP BY Symbol

-- Price per stock per hour
SELECT p.Symbol,
			 p.PostHour,
			 ROUND(MIN(c.AskPrice),2) MinAskPrice,
			 ROUND(AVG(c.AskPrice),2) AvgAskPrice,
			 ROUND(MAX(c.AskPrice),2) MaxAskPrice,
			 ROUND(AVG(p.SentimentScore),3) AvgSentiment,
			 ROUND(AVG(p.WeightedScore),2) WeightedScore,
			 SUM(p.PostCount) PostCount,
			 SUM(p.PositiveCount) Positives,
			 SUM(p.NegativeCount) Negatives,
			 SUM(p.NoiseCount) Noise
FROM (
	SELECT r.Symbol,
	       FORMAT_TIMESTAMP("%F-%k", r.PostedOn)
				 PostHour, r.SentimentScore,
				 -- Using Friend Count as a weight multiplier for significant sentiment, else 0
			   (1 + a.FriendCount / 1000) * CASE
				 			-- if strong then full credit for score
				 			WHEN ABS(r.SentimentScore) > 0.6 THEN r.SentimentScore
							-- if poor give them 50% credit for the sentiment score
							WHEN ABS(r.SentimentScore) >= 0.3 AND ABS(r.SentimentScore) <= 0.6  THEN r.SentimentScore * 0.5
							-- else it is noise to ignore
							ELSE 0 END WeightedScore,
				 1 PostCount,
				 -- for positic/negatives using 0.3 significant sentiment filter, else noise
				 CASE WHEN r.SentimentScore >= 0.3 THEN 1 ELSE 0 END PositiveCount,
				 CASE WHEN r.SentimentScore <= -0.3 THEN 1 ELSE 0 END NegativeCount,
				 CASE WHEN r.SentimentScore > -0.3 AND r.SentimentScore < 0.3 THEN 1 ELSE 0 END NoiseCount
	FROM Posts r
	JOIN Authors a ON r.Username = a.Username
) p
LEFT JOIN (
    SELECT Symbol, FORMAT_TIMESTAMP("%F-%k", SampleOn) AskHour, AskPrice
		FROM Prices
) c ON p.Symbol = c.Symbol AND p.PostHour = c.AskHour
GROUP BY p.Symbol, p.PostHour
ORDER BY p.Symbol, p.PostHour DESC

-- Same as above but for a given stock
SELECT p.Symbol,
			 p.PostHour,
			 ROUND(MIN(c.AskPrice),2) MinAskPrice,
			 ROUND(AVG(c.AskPrice),2) AvgAskPrice,
			 ROUND(MAX(c.AskPrice),2) MaxAskPrice,
			 ROUND(AVG(p.SentimentScore),3) AvgSentiment,
			 ROUND(AVG(p.WeightedScore),2) WeightedScore,
			 SUM(p.PostCount) PostCount,
			 SUM(p.PositiveCount) Positives,
			 SUM(p.NegativeCount) Negatives,
			 SUM(p.NoiseCount) Noise
FROM (
	SELECT r.Symbol,
	       FORMAT_TIMESTAMP("%F-%k", r.PostedOn)
				 PostHour, r.SentimentScore,
				 -- Using Friend Count as a weight multiplier for significant sentiment, else 0
			   (1 + a.FriendCount / 1000) * CASE
				 			-- if strong then full credit for score
				 			WHEN ABS(r.SentimentScore) > 0.6 THEN r.SentimentScore
							-- if poor give them 50% credit for the sentiment score
							WHEN ABS(r.SentimentScore) >= 0.3 AND ABS(r.SentimentScore) <= 0.6  THEN r.SentimentScore * 0.5
							-- else it is noise to ignore
							ELSE 0 END WeightedScore,
				 1 PostCount,
				 -- for positic/negatives using 0.3 significant sentiment filter, else noise
				 CASE WHEN r.SentimentScore >= 0.3 THEN 1 ELSE 0 END PositiveCount,
				 CASE WHEN r.SentimentScore <= -0.3 THEN 1 ELSE 0 END NegativeCount,
				 CASE WHEN r.SentimentScore > -0.3 AND r.SentimentScore < 0.3 THEN 1 ELSE 0 END NoiseCount
	FROM Posts r
	JOIN Authors a ON r.Username = a.Username
	WHERE r.Symbol = 'MSFT'
) p
LEFT JOIN (
    SELECT Symbol, FORMAT_TIMESTAMP("%F-%k", SampleOn) AskHour, AskPrice
		FROM Prices
) c ON p.Symbol = c.Symbol AND p.PostHour = c.AskHour
GROUP BY p.Symbol, p.PostHour
ORDER BY p.Symbol, p.PostHour DESC


-- Select hourly price
SELECT p.Symbol, p.AskHour, ROUND(avg(p.AskPrice),2) as AvgAskPrice, count(*) RecordCount
FROM (
  SELECT Symbol, FORMAT_TIMESTAMP("%F-%k", SampleOn) AskHour, AskPrice
  FROM Prices
) p
GROUP BY p.Symbol, p.AskHour

-- Select posts sentiment by hour
SELECT p.Symbol, p.PostHour, ROUND(AVG(p.SentimentScore),2) as AvgSentiment, count(*) RecordCount
FROM (
  SELECT Symbol, FORMAT_TIMESTAMP("%F-%k", PostedOn) PostHour, SentimentScore
  FROM Posts
) p
GROUP BY p.Symbol, p.PostHour

-- Most recent negative tweets
select p.PostID, FORMAT_TIMESTAMP("%F-%k", p.PostedOn) PostHour, p.SentimentScore, p.Content
FROM Posts p
JOIN Authors a on p.Username = a.Username
WHERE p.SentimentScore < 0 and a.Source = 'Twitter'
ORDER BY p.PostedOn DESC

-- One massive query for negative twitter posts by author's post count
SELECT a.PostCount, p.SentimentScore, p.Content
FROM Authors a
JOIN Posts p ON a.Username = p.Username
JOIN Prices r ON p.Symbol = r.Symbol
WHERE p.Symbol = 'AAPL'
AND a.Source = 'Twitter'
AND p.SentimentScore < 0
ORDER BY a.PostCount DESC
