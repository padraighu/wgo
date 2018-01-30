package news

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