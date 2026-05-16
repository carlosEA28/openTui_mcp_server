package helpers

import (
	"bytes"
	"log"

	"github.com/PuerkitoBio/goquery"
)

// Parser parses the raw HTML fetched from the openTUI documentation and extracts the relevant information to be used in the TUI.
func Parser(html []byte) {

	//Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	//Find the relevant information using CSS selectors
	doc.Find("main").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find("h1").First().Text()
		subtitle := selection.Find("h2").First().Text()
		link, _ := selection.Find("a").First().Attr("href")
		code := selection.Find("pre").First().Text()
		content := selection.Find("content").First().Text()
		article := selection.Find("article").First().Text()

		if title != "" {
			log.Println("title:", title)
		}
		if subtitle != "" {
			log.Println("subtitle:", subtitle)
		}
		if link != "" {
			log.Println("link:", link)
		}
		if code != "" {
			log.Println("code:", code)
		}
		if content != "" {
			log.Println("content:", content)
		}
		if article != "" {
			log.Println("article:", article)
		}
	})

}
