package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/mmcdole/gofeed"
)

type FeedItem struct {
	Title       string
	Link        string
	Description string
}

type SourceFeed struct {
	Host    string
	Favicon string
	Items   []FeedItem
}

var cachedFeeds []SourceFeed
var lastFetched time.Time

func getFeeds(sources []string) []SourceFeed {
	elapsed := time.Since(lastFetched)
	if elapsed < 5*time.Minute {
		log.Printf("It has been %s since last fetching so returning cached feeds.", elapsed.Round(time.Second))
		return cachedFeeds
	}

	var wg sync.WaitGroup
	feedParser := gofeed.NewParser()
	sourceFeeds := make([]SourceFeed, 0)

	for _, source := range sources {
		wg.Add(1)
		go func(src string) {
			defer wg.Done()

			feed, err := feedParser.ParseURL(src)
			if err != nil {
				fmt.Printf("error fetching feed: %v\n", err)
				return
			}

			parsedUrl, err := url.Parse(src)
			if err != nil {
				fmt.Printf("error parsing url: %v\n", err)
				return
			}

			hostname := parsedUrl.Hostname()

			faviconUrl := "https://" + hostname + "/favicon.ico"
			resp, err := http.Get(faviconUrl)
			if err != nil {
				fmt.Printf("error fetching favicon: %v\n", err)
				return
			}
			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("error reading favicon: %v\n", err)
				return
			}

			faviconBase64 := base64.StdEncoding.EncodeToString(data)

			items := make([]FeedItem, 0)
			for _, item := range feed.Items {
				items = append(items, FeedItem{Title: item.Title, Link: item.Link, Description: item.Description})
			}

			sourceFeeds = append(sourceFeeds, SourceFeed{Host: hostname, Favicon: faviconBase64, Items: items})
		}(source)
	}
	wg.Wait()

	cachedFeeds = sourceFeeds
	lastFetched = time.Now()

	return sourceFeeds
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	sources := []string{
		"https://www.ntv.com.tr/son-dakika.rss",
		"https://www.sozcu.com.tr/rss/tum-haberler.xml",
		"https://feeds.bbci.co.uk/turkce/rss.xml",
		"https://www.trthaber.com/manset_articles.rss",
		"https://www.mynet.com/haber/rss/sondakika",
		// add more RSS feed sources here...
	}

	feeds := getFeeds(sources)

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Content-Type", "application/json")

	log.Printf(strconv.Itoa(len(sources)) + " source feeds retrieved!")

	json.NewEncoder(w).Encode(feeds)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/feeds", feedHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./rss-ui/dist/rss-index/"))))
	log.Println("Server listening on http://localhost:8080")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))
}
