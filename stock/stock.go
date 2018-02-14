package stock

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os"
        "encoding/csv"
        "path/filepath"
        "strings"
        "bufio"
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

func GetStocks(symbols []string, c chan<- string) {
        channels := make([]chan string, len(symbols))
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

        c <- fmt.Sprintf("%s - %.2f, %.2f%% %s\n", dat.Symbol, dat.IexRealtimePrice, dat.ChangePercent*100, dat.CompanyName)
}

func ReadStockListConfig() []string {
        wd, err := os.Getwd()
        if err != nil {
                log.Fatal(err)
        }

        f, err := os.Open(filepath.Join(wd, "/stock/stocklist.csv"))
        if err != nil {
                log.Fatal(err)
        }
        defer f.Close()

        reader := csv.NewReader(f)
        table, err := reader.ReadAll()
        if err != nil {
                log.Fatal(err)
        }

        rows := make([]string, len(table))
        for i, row := range table {
                rows[i] = row[0]
        }
        return rows
}

func UpdateStockList(rows []string) {
        stocksStr := strings.Join(rows, ", ")

        // TODO restructure source code 
        newData := fmt.Sprintf(`package stock

const stockList := {%s}`, stocksStr)
        
        wd, err := os.Getwd()
        if err != nil {
                log.Fatal(err)
        }
        
        f, err := os.Create(filepath.Join(wd, "/stock/stocklist.go"))
        if err != nil {
                log.Fatal(err)
        }
        defer f.Close()

        w := bufio.NewWriter(f)
        nWrite, err := w.WriteString(newData)

        if err != nil {
                log.Fatal(err)
        }
        fmt.Println(nWrite)

        w.Flush()
}