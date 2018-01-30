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