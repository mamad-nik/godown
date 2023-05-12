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
	var localSize, actualSize int64
	if _, err := os.Stat(myfilepath); !os.IsNotExist(err) {
		// File exists, compare sizes
		actualSize = getactualfilesize(link)
		localSize = getlocalfilesize(myfilepath)
		if actualSize == localSize {
			fmt.Printf("Skipping %s, file already downloaded\n", filepath.Base(myfilepath))
			return
		} else {
			fmt.Printf("Overwriting %s, actual size: %d bytes, local size: %d bytes\n", filepath.Base(myfilepath), actualSize, localSize)
		}
	}
	/*
	 */
	c := http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	errcheck(err)

	s := fmt.Sprintf("bytes 0-%d/%d", localSize, actualSize)
	fmt.Println(s)
	req.Header.Add("Content-Range", s)

	r := fmt.Sprintf("bytes=%d-%d", localSize, actualSize-1)
	fmt.Println(r)
	req.Header.Add("Range", r)

	err = os.MkdirAll(filepath.Dir(myfilepath), 0755)
	errcheck(err)
	file, err := os.Create(myfilepath)
	errcheck(err)
	resp, err := c.Do(req)
	errcheck(err)
	defer resp.Body.Close()
	size, err := io.Copy(file, resp.Body)
	errcheck(err)
	fmt.Printf("Downloaded file: %s with size %d\n", filepath.Base(myfilepath), size)

	errcheck(err)
	fmt.Println(resp.Status)
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
	fmt.Print("Enter the file path to save the downloaded file: ")

	myfilepath, err := reader.ReadString('\n')
	errcheck(err)
	myfilepath = strings.TrimSpace(myfilepath)
	return theurl, myfilepath
}

func main() {
	var links []string

	theurl, myfilepath := "http://znucomputer.ir/HTML/Semester6/artificial_intelligence.html", "/home/mamad/Downloads"

	getLinks(&links, theurl)

	filename := getfilenames(theurl)
	for i, v := range links {
		fmt.Println(i, v)
		name := fmt.Sprintf("%s/%s/%s%d", myfilepath, filename, filename, i)
		downloadFile(v, name)
	}
}
