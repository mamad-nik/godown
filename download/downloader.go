package download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gocolly/colly"
)

func Errcheck(err error) {

	if err != nil {
		log.Fatal(err)
	}

}
func Getfilenames(theurl string) string {
	regpat := `/([^/]+)\.html$`
	re := regexp.MustCompile(regpat)
	match := re.FindStringSubmatch(theurl)
	return match[1]

}
func GetLinks(links *[]string, theurl string) {
	c := colly.NewCollector()

	c.OnHTML("[data-href]", func(h *colly.HTMLElement) {
		*links = append(*links, h.Attr("data-href"))
		if len(*links) == 0 {
			log.Fatal("there is nothing to download")
		}
	})

	c.Visit(theurl)

}

func DownloadFile(link, myfilepath string) {
	var localSize int64
	actualSize := Getactualfilesize(link)

	if _, err := os.Stat(myfilepath); !os.IsNotExist(err) {
		localSize = Getlocalfilesize(myfilepath)

		if actualSize == localSize {
			fmt.Printf("Skipping %s, file already downloaded\n", filepath.Base(myfilepath))
			return
		} else {
			fmt.Printf("continuing %s, %d MB remaining\n", filepath.Base(myfilepath), (actualSize-localSize)/1048576)
		}

	} else {
		fmt.Printf("Downloading %s with size of %d\n", filepath.Base(myfilepath), actualSize/1048576)
	}
	since := time.Now()
	c := http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	Errcheck(err)

	s := fmt.Sprintf("bytes 0-%d/%d", localSize, actualSize)
	req.Header.Add("Content-Range", s)

	r := fmt.Sprintf("bytes=%d-%d", localSize, actualSize-1)
	req.Header.Add("Range", r)

	err = os.MkdirAll(filepath.Dir(myfilepath), 0755)
	Errcheck(err)
	file, err := os.OpenFile(myfilepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	Errcheck(err)
	resp, err := c.Do(req)
	Errcheck(err)
	defer resp.Body.Close()
	size, err := io.Copy(file, resp.Body)
	Errcheck(err)
	t := time.Since(since)
	fmt.Printf("\n\tDownloaded file: %s with size %d in %v\n\n", filepath.Base(myfilepath), size/1048576, t)

	Errcheck(err)
}

func Getactualfilesize(url string) int64 {
	resp, err := http.Head(url)
	Errcheck(err)
	defer resp.Body.Close()

	size := resp.ContentLength
	Errcheck(err)
	return size
}

func Getlocalfilesize(filepath string) int64 {
	info, err := os.Stat(filepath)
	Errcheck(err)
	return info.Size()
}
