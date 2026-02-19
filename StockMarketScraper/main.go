package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Stock struct {
	company, price, change string
}

func main() {

	ticker := []string{"RVNL.NS",
		"NHPC.NS",
		"TATASTEEL.NS",
		"MSFT",
		"AMZN",
		"SUZLON.NS",
		"GOLDBEES.BO",
		"AAPL",
		"MRF.NS"}

	stocks := []Stock{}

	c := colly.NewCollector()

	//prints when visiting an URL
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla 5.0")
		fmt.Println("visiting", r.URL)
	})

	//error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	//

	c.OnHTML("div#quote-header-info", func(e *colly.HTMLElement) {
		stock := Stock{}
		stock.company = e.ChildText("h1")
		fmt.Println("Company:", stock.company)
		stock.price = e.ChildText("fin-streamer[data-field='regularMarketPrice']")
		fmt.Println("Price:", stock.price)
		stock.change = e.ChildText("fin-streamer[data-field='regularMarketChangePercent']")
		fmt.Println("Change:", stock.change)
		stocks = append(stocks, stock)
	})

	for _, t := range ticker {
		c.Visit("https://finance.yahoo.com/quote/" + t + "/")
	}
	c.Wait()
	fmt.Println(stocks)

	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatalln("csv not created ", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file) //csv new package to work with csv files

	headers := []string{
		"company name",
		"current price",
		"change in price",
	}
	writer.Write(headers) // To write the above heading

	for _, stock := range stocks {
		record := []string{
			stock.company,
			stock.price,
			stock.change,
		}

		writer.Write(record) //write the scraped parameters

	}

	defer writer.Flush()

}
