package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"github.com/lib/pq"
)

var db *sql.DB

type StockNews struct {
	ID             int
	Context        []string
	Summary        string
	News           []string
	Tickers        []string
	DatePublished  time.Time 
	Title          string
	Link           string
}

func main() {
	connStr := "postgres://postgres:password123@localhost:5432/stock_stream?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, context, summary, news, tickers, date_published, title, link FROM cleansed_articles`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var stockNewsList []StockNews

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
			log.Fatal(err)
		}
		stockNewsList = append(stockNewsList, news)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, news := range stockNewsList {
		fmt.Printf("ID: %d, Title: %s, Date Published: %s\n", news.ID, news.Title, news.DatePublished)
	}
}

