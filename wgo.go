package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "time"
)

type NewsAPIResponse struct {
        Status       string `json:"status"`
        TotalResults int    `json:"totalResults"`
        Articles     []struct {
                Source struct {
                        ID   string `json:"id"`
                        Name string `json:"name"`
                } `json:"source"`
                Author      string    `json:"author"`
                Title       string    `json:"title"`
                Description string    `json:"description"`
                URL         string    `json:"url"`
                URLToImage  string    `json:"urlToImage"`
                PublishedAt time.Time `json:"publishedAt"`
        } `json:"articles"`
}

type StockAPIResponse struct {
        Symbol           string  `json:"symbol"`
        CompanyName      string  `json:"companyName"`
        PrimaryExchange  string  `json:"primaryExchange"`
        Sector           string  `json:"sector"`
        CalculationPrice string  `json:"calculationPrice"`
        Open             float64 `json:"open"`
        OpenTime         int64   `json:"openTime"`
        Close            float64 `json:"close"`
        CloseTime        int64   `json:"closeTime"`
        LatestPrice      float64 `json:"latestPrice"`
        LatestSource     string  `json:"latestSource"`
        LatestTime       string  `json:"latestTime"`
        LatestUpdate     int64   `json:"latestUpdate"`
        LatestVolume     int     `json:"latestVolume"`
        IexRealtimePrice float64 `json:"iexRealtimePrice"`
        IexRealtimeSize  int     `json:"iexRealtimeSize"`
        IexLastUpdated   int64   `json:"iexLastUpdated"`
        DelayedPrice     float64 `json:"delayedPrice"`
        DelayedPriceTime int64   `json:"delayedPriceTime"`
        PreviousClose    float64 `json:"previousClose"`
        Change           float64 `json:"change"`
        ChangePercent    float64 `json:"changePercent"`
        IexMarketPercent float64 `json:"iexMarketPercent"`
        IexVolume        int     `json:"iexVolume"`
        AvgTotalVolume   int     `json:"avgTotalVolume"`
        IexBidPrice      float64 `json:"iexBidPrice"`
        IexBidSize       int     `json:"iexBidSize"`
        IexAskPrice      float64 `json:"iexAskPrice"`
        IexAskSize       int     `json:"iexAskSize"`
        MarketCap        int64   `json:"marketCap"`
        PeRatio          float64 `json:"peRatio"`
        Week52High       float64 `json:"week52High"`
        Week52Low        float64 `json:"week52Low"`
        YtdChange        float64 `json:"ytdChange"`
}

type WeatherAPIResponse struct {
        Coord struct {
                Lon float64 `json:"lon"`
                Lat float64 `json:"lat"`
        } `json:"coord"`
        Weather []struct {
                ID          int    `json:"id"`
                Main        string `json:"main"`
                Description string `json:"description"`
                Icon        string `json:"icon"`
        } `json:"weather"`
        Base string `json:"base"`
        Main struct {
                Temp     float64 `json:"temp"`
                Pressure float64 `json:"pressure"`
                Humidity float64 `json:"humidity"`
                TempMin  float64 `json:"temp_min"`
                TempMax  float64 `json:"temp_max"`
        } `json:"main"`
        Visibility int `json:"visibility"`
        Wind       struct {
                Speed float64 `json:"speed"`
                Deg   float64 `json:"deg"`
        } `json:"wind"`
        Clouds struct {
                All int `json:"all"`
        } `json:"clouds"`
        Dt  int `json:"dt"`
        Sys struct {
                Type    int     `json:"type"`
                ID      int     `json:"id"`
                Message float64 `json:"message"`
                Country string  `json:"country"`
                Sunrise int     `json:"sunrise"`
                Sunset  int     `json:"sunset"`
        } `json:"sys"`
        ID   int    `json:"id"`
        Name string `json:"name"`
        Cod  int    `json:"cod"`
}

func GetNews(c chan<- string) {
        resp, err := http.Get("https://newsapi.org/v2/top-headlines?country=us&apiKey=58c28f3f027e47659e2e1815a5604ac2&pageSize=10")
        if err != nil {
                log.Fatal(err)
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Fatal(err)
        }

        var dat NewsAPIResponse
        if err := json.Unmarshal(body, &dat); err != nil {
                log.Fatal(err)
        }
        articles := dat.Articles

        var result string
        for _, article := range articles {
                result += fmt.Sprintf("%s - %s\n", article.Title, article.Source.Name)
        }
        c <- result
}

func GetStocks(c chan<- string) {
        symbols := [3]string{"spy", "dia", "iwm"}
        var channels [3]chan string
        for i, symbol := range symbols {
                channels[i] = make(chan string)
                go GetStock(symbol, channels[i])
        }
        var result string
        for _, channel := range channels {
                result += <-channel
        }
        c <- result
}

func GetStock(symbol string, c chan<- string) {
        resp, err := http.Get(fmt.Sprintf("https://api.iextrading.com/1.0/stock/%s/quote", symbol))
        if err != nil {
                log.Fatal(err)
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Fatal(err)
        }

        var dat StockAPIResponse
        if err := json.Unmarshal(body, &dat); err != nil {
                log.Fatal(err)
        }

        c <- fmt.Sprintf("%s - %.2f, %.2f%%\n", dat.Symbol, dat.IexRealtimePrice, dat.ChangePercent*100)
}

func GetWeather(c chan<- string) {
        resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?zip=77005,us&appid=3cf9bbdefc803717141614702f1f1658&units=metric")
        if err != nil {
                log.Fatal(err)
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Fatal(err)
        }

        var dat WeatherAPIResponse
        if err := json.Unmarshal(body, &dat); err != nil {
                log.Fatal(err)
        }
        c <- fmt.Sprintf("%.2f Celsius %s\n", dat.Main.Temp, dat.Weather[0].Main)
}

type AggregatedInfo struct {
        Weather string
        Stocks  string
        News    string
}

func main() {
        cNews, cStocks, cWeather := make(chan string), make(chan string), make(chan string)
        var info AggregatedInfo

        go GetNews(cNews)
        go GetStocks(cStocks)
        go GetWeather(cWeather)

        info.Weather, info.Stocks, info.News = <-cWeather, <-cStocks, <-cNews
        fmt.Printf("%s\n%s\n%s", info.Weather, info.Stocks, info.News)
}