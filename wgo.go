package main

import (
        "fmt"
        "wgo/weather"
        "wgo/news"
        "wgo/stock"
)

type AggregatedInfo struct {
        Weather string
        Stocks  string
        News    string
}

func main() {
        cNews, cStocks, cWeather := make(chan string), make(chan string), make(chan string)
        var info AggregatedInfo

        go news.GetNews(cNews)
        go stock.GetStocks(cStocks)
        go weather.GetWeather(cWeather)

        info.Weather, info.Stocks, info.News = <-cWeather, <-cStocks, <-cNews
        fmt.Printf("%s\n%s\n%s", info.Weather, info.Stocks, info.News)
}