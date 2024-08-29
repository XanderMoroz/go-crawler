package main

import (
	"log"

	"github.com/gocolly/colly/v2"
)

type product struct {
	Title  string `json:"title"`
	Link   string `json:"url"`
	Price  string `json:"price"`
	Status string `json:"in_stock"`
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("japonica.kz"),
	)

	var myProducts []product

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)

	})

	// // fmt.Println("Ghbdtn")
	// c.OnHTML("body", func(h *colly.HTMLElement) {
	// 	log.Println("Found an body tag!")
	// })

	c.OnHTML("div.catalog-block__inner", func(h *colly.HTMLElement) {
		// log.Println("Found an div tag!")
		item := product{
			Title:  h.ChildText("a.dark_link"),
			Link:   "https://japonica.kz" + h.ChildAttr("a.image-list__link", "href"),
			Price:  h.ChildText("span.price__new-val"),
			Status: h.ChildText("span.status-icon"),
		}
		log.Println("Продукт", item)
		myProducts = append(myProducts, item)
	})

	// c.OnError(func(_ *colly.Response, err error) {
	// 	log.Println("Something went wrong:", err)
	// })

	c.Visit("https://japonica.kz/catalog/yaponskaya-apteka/assortiment/?PAGEN_1=1")

	log.Println("Извлечено", len(myProducts))

	// content, err := json.Marshal(myProducts)

	// if err != nil {
	// 	log.Println(err.Error())
	// }

	// os.WriteFile("products.json", content, 066)
	// c.Visit("https://japonica.kz/catalog/yaponskaya-apteka/assortiment/")
}
