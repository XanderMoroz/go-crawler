package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

type Product struct {
	Title  string `json:"title"`
	Link   string `json:"url"`
	Price  string `json:"price"`
	Status string `json:"in_stock"`
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("japonica.kz"),
	)

	var myProducts []Product

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)

	})

	// // fmt.Println("Ghbdtn")
	// c.OnHTML("body", func(h *colly.HTMLElement) {
	// 	log.Println("Found an body tag!")
	// })

	c.OnHTML("div.catalog-block__inner", func(h *colly.HTMLElement) {
		// log.Println("Found an body tag!")
		// log.Println(h.ChildText("a.dark_link"))
		// log.Println(h.ChildText("span.price__new-val"))
		// log.Println(h.ChildText("span.status-icon"))
		// log.Println(h.ChildAttr("a.dark_link", "href"))
		item := Product{
			Title:  h.ChildText("a.dark_link"),
			Link:   "https://japonica.kz" + h.ChildAttr("a.dark_link", "href"),
			Price:  h.ChildText("span.price__new-val"),
			Status: h.ChildText("span.status-icon"),
		}
		log.Println("Продукт", item)
		myProducts = append(myProducts, item)
	})

	// c.OnHTML("a.arrows-pagination__next", func(h *colly.HTMLElement) {
	// 	next_page := h.Request.AbsoluteURL(h.Attr("href"))
	// 	log.Println(h.Request.AbsoluteURL(h.Attr("href")))
	// 	c.Visit(next_page)

	// })

	c.Visit("https://japonica.kz/catalog/podarochnye-sertifikaty/")

	log.Println("Извлечено", len(myProducts))

	content, err := json.Marshal(myProducts)

	if err != nil {
		log.Println(err.Error())
	}

	os.WriteFile("products.json", content, 0644)

}
