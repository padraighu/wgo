package cryptocurrency

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "time"
)

type CoindeskAPIResponse struct {
	Time struct {
		Updated    string    `json:"updated"`
		UpdatedISO time.Time `json:"updatedISO"`
		Updateduk  string    `json:"updateduk"`
	} `json:"time"`
	Disclaimer string `json:"disclaimer"`
	ChartName  string `json:"chartName"`
	Bpi        struct {
		USD struct {
			Code        string  `json:"code"`
			Symbol      string  `json:"symbol"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"USD"`
		GBP struct {
			Code        string  `json:"code"`
			Symbol      string  `json:"symbol"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"GBP"`
		EUR struct {
			Code        string  `json:"code"`
			Symbol      string  `json:"symbol"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"EUR"`
	} `json:"bpi"`
}

func GetBitcoinPrice(c chan<- string) {
        resp, err := http.Get("https://api.coindesk.com/v1/bpi/currentprice.json")
        if err != nil {
                log.Fatal(err)
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Fatal(err)
        }

        var dat CoindeskAPIResponse
        if err := json.Unmarshal(body, &dat); err != nil {
                log.Fatal(err)
        }

        c <- fmt.Sprintf("Bitcoin - $%.2f\n", dat.Bpi.USD.RateFloat) 
}