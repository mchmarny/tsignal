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
