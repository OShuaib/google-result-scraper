package main

import(
	"fmt"
	"time"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SearchResult struct {
	ResultRank int
	ResultURL string
	ResultTitle string
	ResultDesc	string
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

func randomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func buildGoogleUrls(searchTerm, countryCode, languageCode string, pages, count int)([]string, error){
	toScrape := []string{}
	searchTerm = strings.Trim(searchTerm," ")
	searchTerm = strings.Replace(searchTerm," ", "+",-1)
	if googleBase, found := googleDomains[countryCode]; found {
		for i := 0; i < pages; i++ {
			start := i * count
			scrapeURL := fmt.Sprintf("%s%s&num=%d&hl=%s&start=%d&filter=0",googleBase,searchTerm, count, languageCode, start)
			toScrape = append(toScrape, scrapeURL)
		}
	}else {
		err := fmt.Errorf("country (%s) is currently not supported", countryCode)
		return nil, err
	}
	return toScrape,nil

}

func googleResultParsing(response *http.Response, rank int) ([]SearchResult, error){
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil,err
	}
	results := []SearchResult{}
	sel := doc.Find("div.g")
	rank ++ 
	for i := range sel.Nodes{
		item:= sel.Eq(i)
		linkTag:= item.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h3.r")
		descTag := item.Find("span.st")
		desc := descTag.Text()
		title := titleTag.Text()
		link = strings.Trim(link, " ")

		if link != "" && link != "#" && !strings.HasPrefix(link,"/") {
			result := SearchResult {
				rank,
				link,
				title,
				desc,
			}
			results = append(results, result)
			rank ++
		}
	}
	return results, err
}

