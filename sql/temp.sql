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
