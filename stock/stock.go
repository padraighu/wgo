package stock

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
)

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