package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gocolly/colly"
)

func errcheck(err error) {
	if err != nil {
		log.Fatal(err)
	}

}
func getfilenames(theurl string) string {
	regpat := `/([^/]+)\.html$`
	re := regexp.MustCompile(regpat)
	match := re.FindStringSubmatch(theurl)
	return match[1]

}
func getLinks(links *[]string, theurl string) {
	c := colly.NewCollector()

	c.OnHTML("[data-href]", func(h *colly.HTMLElement) {
		*links = append(*links, h.Attr("data-href"))
	})

	c.Visit(theurl)

}

func downloadFile(filename, link string) {
	file, err := os.Create(filename)
	errcheck(err)
	resp, err := http.Get(link)
	errcheck(err)
	defer resp.Body.Close()
	size, err := io.Copy(file, resp.Body)
	errcheck(err)
	fmt.Printf("Downloaded file: %s with size %d\n", filename, size)

}

func main() {
	var links []string
	theurl := "http://znucomputer.ir/HTML/Semester6/artificial_intelligence.html"

	getLinks(&links, theurl)
	filename := getfilenames(theurl)
	for i, v := range links {
		name := fmt.Sprintf("%s%d", filename, i)
		downloadFile(name, v)
	}

}
