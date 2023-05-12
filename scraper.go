package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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
	var localSize int64
	actualSize := getactualfilesize(link)

	if _, err := os.Stat(myfilepath); !os.IsNotExist(err) {
		localSize = getlocalfilesize(myfilepath)

		if actualSize == localSize {
			fmt.Printf("Skipping %s, file already downloaded\n", filepath.Base(myfilepath))
			return
		} else {
			fmt.Printf("continuing %s, %d MB remaining\n", filepath.Base(myfilepath), (actualSize-localSize)/1048576)
		}

	} else {
		fmt.Printf("Downloading %s with size of %d\n", filepath.Base(myfilepath), actualSize/1048576)
	}

	c := http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	errcheck(err)

	s := fmt.Sprintf("bytes 0-%d/%d", localSize, actualSize)
	req.Header.Add("Content-Range", s)

	r := fmt.Sprintf("bytes=%d-%d", localSize, actualSize-1)
	req.Header.Add("Range", r)
	start := time.Now()
	err = os.MkdirAll(filepath.Dir(myfilepath), 0755)
	errcheck(err)
	file, err := os.OpenFile(myfilepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	errcheck(err)
	resp, err := c.Do(req)
	errcheck(err)
	defer resp.Body.Close()
	size, err := io.Copy(file, resp.Body)
	errcheck(err)
	elapsed := time.Since(start)
	fmt.Printf("\n\tDownloaded file: %s with size %d in %v\n\n", filepath.Base(myfilepath), size/1048576, elapsed)

	errcheck(err)
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

// "http://znucomputer.ir/HTML/Semester6/artificial_intelligence.html"
// "/home/mamad/Downloads"

func userInput() (string, string) {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("enter URL of the page you want to download from: ")
	theurl, err := reader.ReadString('\n')
	errcheck(err)

	theurl = strings.TrimSpace(theurl)
	fmt.Print("Enter the file absolute path to save the downloaded file: ")

	myfilepath, err := reader.ReadString('\n')
	errcheck(err)
	myfilepath = strings.TrimSpace(myfilepath)
	return theurl, myfilepath
}

func main() {
	var links []string

	theurl, myfilepath := userInput()

	getLinks(&links, theurl)

	filename := getfilenames(theurl)
	for i, v := range links {
		fmt.Printf("%s: ", v)
		name := fmt.Sprintf("%s/%s/%s%d", myfilepath, filename, filename, i)
		downloadFile(v, name)
	}
}
