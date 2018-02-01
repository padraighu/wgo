package main

import (
        "fmt"
        "wgo/weather"
        "wgo/news"
        "wgo/stock"
        "wgo/cryptocurrency"
)

type AggregatedInfo struct {
        Weather string
        Stocks  string
        News    string
        Crypto  string
}

func main() {
        cNews, cStocks, cWeather, cCrypto := make(chan string), make(chan string), make(chan string), make(chan string)
        var info AggregatedInfo

        go news.GetNews(cNews)
        go stock.GetStocks(cStocks)
        go weather.GetWeather(cWeather)
        go cryptocurrency.GetBitcoinPrice(cCrypto)

        info.Weather, info.Stocks, info.News, info.Crypto = <-cWeather, <-cStocks, <-cNews, <-cCrypto
        fmt.Printf("%s\n%s\n%s\n%s", info.Weather, info.Stocks, info.Crypto, info.News)
}