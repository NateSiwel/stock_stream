package main

import (
	"database/sql"
	"log"
	"github.com/lib/pq"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"strings"
	"fmt"
	"io/ioutil"
	"encoding/json"
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

type NewsItem struct {
    Title   string
    Summary string
}

type StockData struct {
	Datetime map[int]string  `json:"Datetime"`
	AdjClose map[int]float64   `json:"Adj Close"`
	Close    map[int]float64   `json:"Close"`
	High     map[int]float64   `json:"High"`
	Low      map[int]float64   `json:"Low"`
	Open     map[int]float64   `json:"Open"`
	Volume   map[int]float64   `json:"Volume"`
}

type DataPoint struct {
	Label string  `json:"label"`
	Y     float64 `json:"y"`
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

func fetchMarketNews() []NewsItem {
    return []NewsItem{
        {"Market Update 1", "U.S. stock markets had a mixed session, with the S&P 500 up 0.5% and the Dow slightly down by over 20 points. The latest ADP report showed a gain of 150,000 jobs in June, indicating a softening yet stable labor market. Treasury yields fell below 4.4% as economic growth slows, hinting at potential easing in inflation. Technology stocks led gains, benefiting from lower bond yields. Key economic indicators suggest a gradual normalization of the labor market, supporting consumer spending. Investors remain optimistic, anticipating potential Federal Reserve rate cuts later this year."},
    }
}

func main() {
	connStr := "postgres://postgres:password123@localhost:5432/stock_stream?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	resp, err := http.Get("http://localhost:5000/data")
	if err != nil {
	fmt.Println("Error:", err)
	return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var data StockData
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error:", err)
		return
	}

	r := gin.Default()

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		marketNews := fetchMarketNews()
		c.HTML(http.StatusOK, "index.html", gin.H{"marketNews": marketNews,})
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
			retHTML += `<div id="summary-block"><div id="summary-ticker">` + summary.Ticker+ `</div><div id="summary-content">` + summary.Summary+ `</div></div>`
			if index < len(stockSummaries)-1{
				retHTML += `<hr style="visibility:hidden;height:1px;border-width:0;color:gray;background-color:rgb(224,224,224);width:50%;display:flex;align-items:center;">`
			}
		}
		c.String(http.StatusOK, retHTML)
	})

	r.GET("/data", func(c *gin.Context) {
		dataPoints := extractDataPoints(data)
		c.JSON(http.StatusOK, dataPoints)
	})

	r.Run(":8080")

}

func extractDataPoints(data StockData) []DataPoint {
	var dataPoints []DataPoint

	for i := 0; i < len(data.Datetime); i++ {
		datetime := data.Datetime[i]
		closeValue := data.Close[i]

		dataPoints = append(dataPoints, DataPoint{
			Label: datetime,
			Y:     closeValue,
		})
	}

	return dataPoints
}
