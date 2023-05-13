package main

import (
	"fmt"
	"log"
	"os/user"
	"path/filepath"

	"scraper.go/download"
	"scraper.go/page"
)

// "http://znucomputer.ir/HTML/Semester6/artificial_intelligence.html"
// "/home/mamad/Downloads"
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
	var theurl, myfilepath string

	currentUser, err := user.Current()
	download.Errcheck(err)
	downloadsPath := filepath.Join(currentUser.HomeDir, "Downloads/lessons")

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

	myfilepath = downloadsPath
	download.GetLinks(&links, theurl)
	if len(links) == 0 {
		log.Fatal("there is nothing to download")
	}

	filename := download.Getfilenames(theurl)
	for i, v := range links {
		fmt.Printf("%s: ", v)
		name := fmt.Sprintf("%s/%s/%s%d", myfilepath, filename, filename, i)
		download.DownloadFile(v, name)
	}

}
