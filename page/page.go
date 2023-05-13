package page

import (
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func linkgeneratation(base, page string) string {
	regpat := `^([\./])+`
	p := regexp.MustCompile(regpat)
	a := p.FindStringSubmatch(page)
	b := strings.TrimLeft(page, a[0])
	return base + "/" + b
}

func GetPages() *map[string]string {
	class := make(map[string]string)
	link := "http://znucomputer.ir"
	c := colly.NewCollector()
	c.OnHTML(".dropdown-content a", func(h *colly.HTMLElement) {
		class[strings.TrimSpace(h.Text)] = h.Attr("href")
	})
	err := c.Visit(link)
	if err != nil {
		log.Fatal(err)
	}

	for i, v := range class {
		class[i] = linkgeneratation(link, v)
	}
	return &class
}
