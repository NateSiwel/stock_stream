package main

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestFetchStockNews(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "context", "summary", "news", "tickers", "date_published", "title", "link"}
	rows := sqlmock.NewRows(columns).
		AddRow(1, pq.Array([]string{"context1"}), "summary1", pq.Array([]string{"news1"}), pq.Array([]string{"ticker1"}), time.Now(), "title1", "link1")

	tickers := []string{"AAPL", "GOOGL"}

	mock.ExpectQuery(`SELECT id, context, summary, news, tickers, date_published, title, link FROM cleansed_articles WHERE tickers && \$1`).
		WithArgs(pq.Array(tickers)).
		WillReturnRows(rows)

	stockNewsList, err := fetchStockNews(db, tickers)
	assert.NoError(t, err)
	assert.Len(t, stockNewsList, 1)
	assert.Equal(t, 1, stockNewsList[0].ID)
	assert.Equal(t, "title1", stockNewsList[0].Title)
}


