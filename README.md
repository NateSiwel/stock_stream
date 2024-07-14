
# Stock Stream

## TODO
- web scraper that scrapes stock articles and interacts with gpt to provide data to psql database
- go API that interacts w/ postgres to fetch daily AI summarized news data tailored around a users portfolio
  - accepts list of tickers and timeframe
  - API should perform vector search on psql news data for past X timeframe to fetch context for LLM
  - pass news data to llm as context, instructing LLM to output the most important data for user considering their porfolio
  - return output to user - providing useful summary of news tailored for user. 
- Lightweight HTMX frontend

### gpt article cleaner
 > I am going to ask you to return JSON data about an article in the following format w/ the following variables:
 {"context":[
  "Cash flows",
  "Financial stability",
  "Dividend payments",
  "Apple",
  "Microsoft",
  "Alphabet",
  "Cloud revenue",
  "AI",
  "China iPhone",
  "Zacks Rank",
  "Earnings expectations",
  "Bill Gates",
  "Foundation portfolio",
  "Azure cloud",
  "Peloton",
  "Revenue decline",
  "Connected-fitness",
  "Deep-value",
  "AI stocks",
  "Palantir"
],
"summary": "The article discusses the strong financial stability and cash flow generation of three mega-cap companies: Apple, Microsoft, and Alphabet, highlighting their robust performance and growth potential. Apple has seen a recent 15% stock increase due to favorable news about its iPhone shipments in China and AI developments, with consistent dividend growth. Microsoft has posted strong cloud results, leading the market's charge with a 25% increase in share value, driven by AI-related services and cloud revenue growth. Alphabet has experienced a nearly 40% stock increase in 2024, with positive earnings revisions and strong growth prospects. Overall, these companies are well-positioned to weather downturns and offer long-term benefits to investors.",
"news": ["Apple shares have jumped back into favor after a slow start to 2024", "Favorable news concerning Apple's China iPhone shipments has aided performance, with recent commentary surrounding AI also providing tailwinds.", "Apple raised its quarterly payout in its latest quarterly release, reflecting the 12th consecutive year of higher payouts", "Analysts have revised their microsoft earnings expectations for its current fiscal year consistently over the last year, with the $11.77 per share expected up roughly 9% over the last year", etc ..... ],
"tickers":[
    "MSFT",
    "AAPL",
    "GOOGL",
    "PTON",
    "PLTR",
    "NVDA",
    "TSLA"
]}
You will replace the content of the JSON completely with information about the following article:

### API output prompt
> Your job is to give daily market updates customized to a persons portfolio given up-to-date articles. 
Give a quick data-rich 3 sentence market summary per stock, to a person who owns {stocks_string}, given the below up to date data. You shouldn't give ratings, but only report news. Avoid giving data that the user may have already read over 1 day ago.

> Your job is to give daily market updates for the general market based on up-to-date articles. Provide a concise, data-rich six-sentence market summary that covers major indices, key economic indicators, and significant market-moving events or trends. Avoid repeating information from articles that are over one day old and focus on the most current news and data. Return output in {"thinking":"", "daily_market_summary":""} format. You may plan out your article in thinking.
