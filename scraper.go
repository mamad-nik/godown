package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"scraper.go/download"
)

// "http://znucomputer.ir/HTML/Semester6/artificial_intelligence.html"
// "/home/mamad/Downloads"

func userInput() (string, string) {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("enter URL of the page you want to download from: ")
	theurl, err := reader.ReadString('\n')
	download.Errcheck(err)

	theurl = strings.TrimSpace(theurl)
	fmt.Print("Enter the file path to save the downloaded file: ")

	myfilepath, err := reader.ReadString('\n')
	download.Errcheck(err)
	myfilepath = strings.TrimSpace(myfilepath)
	return theurl, myfilepath
}

func main() {
	var links []string

	theurl, myfilepath := "http://znucomputer.ir/HTML/Semester6/artificial_intelligence.html", "/home/mamad/Downloads"

	download.GetLinks(&links, theurl)

	filename := download.Getfilenames(theurl)
	for i, v := range links {
		fmt.Printf("%s: ", v)
		name := fmt.Sprintf("%s/%s/%s%d", myfilepath, filename, filename, i)
		download.DownloadFile(v, name)
	}
}
