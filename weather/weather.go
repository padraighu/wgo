package weather

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
)

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

type IPAPIResponse struct {
        As          string  `json:"as"`
        City        string  `json:"city"`
        Country     string  `json:"country"`
        CountryCode string  `json:"countryCode"`
        Isp         string  `json:"isp"`
        Lat         float64 `json:"lat"`
        Lon         float64 `json:"lon"`
        Org         string  `json:"org"`
        Query       string  `json:"query"`
        Region      string  `json:"region"`
        RegionName  string  `json:"regionName"`
        Status      string  `json:"status"`
        Timezone    string  `json:"timezone"`
        Zip         string  `json:"zip"`
}

func GetZip() string {
        resp, err := http.Get("http://ip-api.com/json")
        if err != nil {
                log.Fatal(err)
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Fatal(err)
        }

        var dat IPAPIResponse
        if err := json.Unmarshal(body, &dat); err != nil {
                log.Fatal(err)
        }

        return dat.Zip
}

func GetWeather(c chan<- string) {
        // Get location via local IP.
        // TODO Consider the case when zip is not available (e.g. non US). Should use city name instead.
        zip := GetZip()

        resp, err := http.Get(fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?zip=%s,us&appid=3cf9bbdefc803717141614702f1f1658&units=metric", zip)) // Metric system FTW
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