package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

type Product struct {
	Title        string            `json:"title"`
	Price        string            `json:"price"`
	Description  string            `json:"description"`
	ProductURL   string            `json:"url"`
	Availability []productInStores `json:"availability"`
}

type productInStores struct {
	StoreAdress string `json:"adress"`
	Status      string `json:"in_stock"`
}

var Products []Product

func main() {
	// Инициализируем новый коллектор
	c := colly.NewCollector(
		colly.AllowedDomains("japonica.kz"),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)

	})

	// Обрабатываем элементы каталога
	c.OnHTML("div.catalog-block__item", func(e *colly.HTMLElement) {
		// log.Println("Обнаружили позицию в каталоге")
		relativeURL := e.ChildAttr("a.dark_link", "href")
		// log.Println(relativeURL)
		courseURL := e.Request.AbsoluteURL(relativeURL)
		// log.Println("Посещаем страницу товара")
		e.Request.Visit(courseURL)
	})

	// Обрабатываем страницу товара
	c.OnHTML("div.catalog-detail__item", func(e *colly.HTMLElement) {

		// Извлекаем URL текущей страницы
		pageURL := e.Request.URL.String()

		// Извлекаем данные о товаре
		title := e.ChildText("h1.font_32")
		price := e.ChildText("span.price__new-val")
		description := e.ChildText("div.content")

		// log.Println(title, price, description, pageURL)
		// price := e.ChildText(".product-price")
		// description := e.ChildText(".product-description")

		// fmt.Printf("Title: %s\nPrice: %s\nDescription: %s\n", title, price, description)
		var InStores []productInStores

		// Обходим список магазинов, в которых товар есть в наличии
		e.ForEach("div.stores-list__item", func(_ int, el *colly.HTMLElement) {
			store := el.ChildText("div.stores-list__item-title")
			stocks := el.ChildText("span.status-icon")
			// log.Println("Обнаружили магазин", el.ChildText("div.stores-list__item-title"))
			// log.Println("Обнаружили статус товара", el.ChildText("span.status-icon"))
			productData := productInStores{
				StoreAdress: store,
				Status:      stocks,
			}

			InStores = append(InStores, productData)
		})

		productData := Product{
			Title:        title,
			Price:        price,
			Description:  description,
			ProductURL:   pageURL,
			Availability: InStores,
		}

		Products = append(Products, productData)
		// log.Println(productData)

	})

	// /catalog/yaponskaya-apteka/assortiment/?PAGEN_1=2
	// /catalog/yaponskaya-apteka/assortiment/?PAGEN_1=3

	// Запускаем скрапинг
	err := c.Visit("https://japonica.kz/catalog/yaponskaya-apteka/assortiment/?PAGEN_1=3")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Извлечено", len(Products))

	// Экспортируем данные в разные форматы файлов
	exportToTXT(Products)
	exportToCSV(Products)
	exportToJSON(Products)

}

func exportToTXT(products []Product) {
	file, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, product := range products {
		_, err = file.WriteString(fmt.Sprintf("%+v\n", product))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func exportToCSV(products []Product) {
	file, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, product := range products {
		err = writer.Write([]string{product.Title, product.Price, product.Description, product.ProductURL})
		if err != nil {
			log.Fatal(err)
		}

		for _, store := range product.Availability {
			err = writer.Write([]string{store.StoreAdress, store.Status})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func exportToJSON(products []Product) {
	file, err := os.Create("output.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(products)
	if err != nil {
		log.Fatal(err)
	}
}

// func main() {
// 	c := colly.NewCollector(
// 		colly.AllowedDomains("japonica.kz"),
// 	)

// 	// Create another collector to scrape course details
// 	// detailCollector := c.Clone()

// 	// var myProducts []Product

// 	c.OnRequest(func(r *colly.Request) {
// 		log.Println("Visiting", r.URL)

// 	})

// 	// // fmt.Println("Ghbdtn")
// 	// c.OnHTML("body", func(h *colly.HTMLElement) {
// 	// 	log.Println("Found an body tag!")
// 	// })

// 	c.OnHTML("div.catalog-detail__top-info", func(h *colly.HTMLElement) {
// 		// log.Println("Found an body tag!")
// 		// log.Println(h.ChildText("a.dark_link"))
// 		// log.Println(h.ChildText("span.price__new-val"))
// 		// log.Println(h.ChildText("span.status-icon"))
// 		// log.Println(h.ChildAttr("a.dark_link", "href"))
// 		// item := Product{
// 		// 	Title:  h.ChildText("h1.font_32"),
// 		// 	Link:   "https://japonica.kz" + h.ChildAttr("a.dark_link", "href"),
// 		// 	Price:  h.ChildText("span.price__new-val"),
// 		// 	Status: h.ChildText("span.status-icon"),
// 		// }
// 		log.Println(h.ChildText("h1.font_32"))
// 		// Activate detailCollector if the link contains "coursera.org/learn"
// 		// courseURL := h.Request.AbsoluteURL(e.Attr("href"))
// 		// if strings.Index(courseURL, "coursera.org/learn") != -1 {
// 		// 	detailCollector.Visit(courseURL)
// 		// }
// 		// myProducts = append(myProducts, item)
// 	})
// 	// c.OnHTML("a.dark_link", func(h *colly.HTMLElement) {
// 	// 	next_page := h.Request.AbsoluteURL(h.Attr("href"))
// 	// 	log.Println(h.Request.AbsoluteURL(h.Attr("href")))
// 	// 	c.Visit(next_page)
// 	// })
// 	// c.OnHTML("a.arrows-pagination__next", func(h *colly.HTMLElement) {
// 	// 	next_page := h.Request.AbsoluteURL(h.Attr("href"))
// 	// 	log.Println(h.Request.AbsoluteURL(h.Attr("href")))
// 	// 	c.Visit(next_page)

// 	// })

// 	c.Visit("https://japonica.kz/catalog/podarochnye-sertifikaty/")

// 	// log.Println("Извлечено", len(myProducts))

// 	// content, err := json.Marshal(myProducts)

// 	// if err != nil {
// 	// 	log.Println(err.Error())
// 	// }

// 	// os.WriteFile("products.json", content, 0644)

// }
