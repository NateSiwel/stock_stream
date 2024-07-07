package main

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/lib/pq"
	"time"
)

type StockNews struct {
	ID             int
	Context        []string
	Summary        string
	News           []string
	Tickers        []string
	DatePublished  string 
	Title          string
	Link           string
}

type StockSummaries struct {
	Ticker	       string
	Summary        string
	Date	       time.Time

}

func fetchStockNews(db *sql.DB, tickers []string) ([]StockNews, error) {
	query := `
		SELECT id, context, summary, news, tickers, date_published, title, link
		FROM cleansed_articles
		WHERE tickers && $1 
	`

	rows, err := db.Query(query, pq.Array(tickers))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stockNews []StockNews

	for rows.Next() {
		var news StockNews
		err := rows.Scan(
			&news.ID,
			pq.Array(&news.Context),
			&news.Summary,
			pq.Array(&news.News),
			pq.Array(&news.Tickers),
			&news.DatePublished,
			&news.Title,
			&news.Link,
		)
		if err != nil {
			return nil, err
		}
		stockNews= append(stockNews, news)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stockNews, nil
}

func fetchStockSummaries(db *sql.DB, tickers []string) ([]StockSummaries, error) {
	query := `
		SELECT ticker, summary, date 
		FROM stock_summaries 
		WHERE ticker = ANY($1);
	`

	rows, err := db.Query(query, pq.Array(tickers))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stockSummaries []StockSummaries

	for rows.Next() {
		var news StockSummaries
		err := rows.Scan(
			&news.Ticker,
			&news.Summary,
			&news.Date,
		)
		if err != nil {
			return nil, err
		}
		stockSummaries = append(stockSummaries, news)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stockSummaries, nil
}

func main() {
	connStr := "postgres://postgres:password123@localhost:5432/stock_stream?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tickers := []string{"AAPL", "PYPL"}

	stockSummaries, err := fetchStockSummaries(db, tickers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User owns:")
	for _, ticker := range tickers {
	fmt.Printf(" %s", ticker)
	}
	fmt.Println()

	for _, summary := range stockSummaries {
		fmt.Printf("Ticker: %s, Summary: %s\n", summary.Ticker, summary.Summary)

	}	
}
