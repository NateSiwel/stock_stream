package main

import (
	"database/sql"
	"log"
	"github.com/lib/pq"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"strings"
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

	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/load-summaries", func(c *gin.Context) {
   	        tickersParam := c.Query("tickers")
	        tickers := strings.Split(tickersParam, ",")

		stockSummaries, err := fetchStockSummaries(db, tickers)
		if err != nil {
			log.Fatal(err)
		}

		retHTML := ""
		for index, summary := range stockSummaries{
			retHTML += `<div id="summary-ticker">` + summary.Ticker+ `</div>`
			retHTML += `<div id="summary-content">` + summary.Summary+ `</div>`
			if index < len(stockSummaries)-1{
				retHTML += `<hr style="height:1px;border-width:0;color:gray;background-color:rgb(224,224,224);width:50%;display:flex;align-items:center;">`
			}
		}
		c.String(http.StatusOK, retHTML)
	})

	r.POST("/send-message", func(c *gin.Context) {
		message := c.PostForm("message")
		messages = append(messages, Message{Sender: "user", Text: message})
		messages = append(messages, Message{Sender: "bot", Text: "Here is the information you requested."})

		messageHTML := `<div class="chat-message user">` + message + `</div>`
		messageHTML += `<div class="chat-message bot">Here is the information you requested.</div>`
		c.String(http.StatusOK, messageHTML)
	})



	r.Run(":8080")

}
