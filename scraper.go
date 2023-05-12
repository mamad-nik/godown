package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func downloadFile(link, myfilepath string) {

	if _, err := os.Stat(myfilepath); !os.IsNotExist(err) {
		// File exists, compare sizes
		actualSize := getactualfilesize(link)
		localSize := getlocalfilesize(myfilepath)
		if actualSize == localSize {
			fmt.Printf("Skipping %s, file already downloaded\n", filepath.Base(myfilepath))
			return
		} else {
			fmt.Printf("Overwriting %s, actual size: %d bytes, local size: %d bytes\n", filepath.Base(myfilepath), actualSize, localSize)
		}
	}

	err := os.MkdirAll(filepath.Dir(myfilepath), 0755)
	errcheck(err)
	file, err := os.Create(myfilepath)
	errcheck(err)
	resp, err := http.Get(link)
	errcheck(err)
	defer resp.Body.Close()
	size, err := io.Copy(file, resp.Body)
	errcheck(err)
	fmt.Printf("Downloaded file: %s with size %d\n", filepath.Base(myfilepath), size)

}

func getactualfilesize(url string) int64 {
	resp, err := http.Head(url)
	errcheck(err)
	defer resp.Body.Close()

	size := resp.ContentLength
	errcheck(err)
	return size
}

func getlocalfilesize(filepath string) int64 {
	info, err := os.Stat(filepath)
	errcheck(err)
	return info.Size()
}

func main() {
	var links []string
	theurl := "http://znucomputer.ir/HTML/Semester6/artificial_intelligence.html"
	myfilepath := "/home/mamad/Downloads"

	getLinks(&links, theurl)

	filename := getfilenames(theurl)
	for i, v := range links {
		name := fmt.Sprintf("%s/%s/%s%d", myfilepath, filename, filename, i)
		downloadFile(v, name)
	}

}
