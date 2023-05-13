package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"scraper.go/download"
	"scraper.go/page"
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
func userInput2(list *[]map[string]string) int {
	var a int
	fmt.Printf("choose the class you want to download\n\n")
	for i, v := range *list {
		for j := range v {
			fmt.Printf("\t%d -> %s\n", i, j)
		}
	}
	fmt.Scanf("%d", &a)
	return a
}
func getList(class *map[string]string, list *[]map[string]string) {
	i := 0
	for k, v := range *class {
		m := make(map[string]string)
		m[k] = v
		*list = append(*list, m)
		i++
	}

}
func main() {
	var list []map[string]string
	var theurl string
	class := page.GetPages()
	getList(class, &list)
	a := userInput2(&list)
	if a >= len(list) {
		log.Fatal("WTF???")
	}
	for i, v := range list[a] {
		fmt.Println(i)
		theurl = v
	}

	var links []string

	myfilepath := "/home/mamad/Downloads"

	download.GetLinks(&links, theurl)

	filename := download.Getfilenames(theurl)
	for i, v := range links {
		fmt.Printf("%s: ", v)
		name := fmt.Sprintf("%s/%s/%s%d", myfilepath, filename, filename, i)
		download.DownloadFile(v, name)
	}

}
