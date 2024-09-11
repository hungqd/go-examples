package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/gocolly/colly/v2"
	"github.com/hungqd/books-crawler/book"
)

func main() {
	books := make(chan *book.Book)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnHTML(`article[class="product_pod"]`, func(e *colly.HTMLElement) {
		// Thumbnail
		thumbnail := e.DOM.Find("img.thumbnail")
		src := thumbnail.AttrOr("src", "")
		thumbnailURL := e.Request.AbsoluteURL(src)

		// Rating
		ratingElem := e.DOM.Find("p.star-rating")
		var rating int = 0
		if ratingElem.HasClass("One") {
			rating = 1
		} else if ratingElem.HasClass("Two") {
			rating = 2
		} else if ratingElem.HasClass("Three") {
			rating = 3
		} else if ratingElem.HasClass("Four") {
			rating = 4
		} else if ratingElem.HasClass("Five") {
			rating = 5
		}

		// Details URL & title
		titleElem := e.DOM.Find("h3 > a")
		detailURL := e.Request.AbsoluteURL(titleElem.AttrOr("href", ""))
		title := titleElem.AttrOr("title", "")

		// Price
		priceElem := e.DOM.Find(".product_price .price_color")
		price := priceElem.Text()

		// Instock
		intockElem := e.DOM.Find(".product_price .instock")
		instockText := strings.TrimSpace(intockElem.Text())
		instock := strings.EqualFold("In stock", instockText)

		books <- &book.Book{
			Thumbnail: thumbnailURL,
			DetailURL: detailURL,
			Rating:    rating,
			Title:     title,
			Price:     price,
			Instock:   instock,
		}
	})

	c.OnHTML(".pager .next a[href]", func(e *colly.HTMLElement) {
		select {
		case <-stop:
			return
		default:
			e.Request.Visit(e.Attr("href"))
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("On error: %v\n", err)
		stop <- syscall.SIGABRT
		close(books)
	})

	var wg sync.WaitGroup

	go func() {
		defer wg.Done()
		for book := range books {
			fmt.Println(book.Title, book.DetailURL)
		}
	}()

	c.Visit("http://books.toscrape.com/index.html")
	close(books)

	wg.Wait()
}
