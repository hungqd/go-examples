package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

var ratingMap = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
}

func getRating(ratingElemClass string) int {
	classes := strings.Split(ratingElemClass, " ")
	if len(classes) < 2 {
		log.Fatalf("Rating element should has more than 1 class (class=%s)\n", ratingElemClass)
	}
	cls := classes[len(classes)-1]
	rating := ratingMap[strings.ToLower(cls)]
	return rating
}

func takeScreenShot(ctx context.Context, description string) {
	var buf []byte
	err := chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf))
	if err != nil {
		log.Printf("Take screenshot error: %v\n", err)
		return
	}
	screenshotDir := filepath.Join(".", "screenshots")
	err = os.MkdirAll(screenshotDir, os.ModePerm)
	if err != nil {
		log.Printf("Create screenshot directory error: %v\n", err)
		return
	}

	outfile := filepath.Join(screenshotDir, time.Now().Format("20060102150405")+"_"+description+".png")
	if err := os.WriteFile(outfile, buf, 0o644); err != nil {
		log.Printf("Save screenshot error: %v\n", err)
	}
}

func main() {
	// allocator, cancel := chromedp.NewRemoteAllocator(
	// 	context.Background(),
	// 	"ws://localhost:9222",
	// )
	// defer cancel()

	// ctx, cancel := chromedp.NewContext(allocator)
	// defer cancel()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var bookNodes []*cdp.Node

	err := chromedp.Run(ctx,
		chromedp.Navigate("http://books.toscrape.com/index.html"),
	)

	if err != nil {
		log.Fatal(err)
	}

	for {
		var url string
		chromedp.Run(ctx,
			chromedp.WaitVisible("article.product_pod", chromedp.ByQueryAll),
			chromedp.Location(&url),
			chromedp.Nodes("article.product_pod", &bookNodes, chromedp.ByQueryAll),
		)

		fmt.Println("===============================================")
		fmt.Printf("PAGE: %s\n", url)

		takeScreenShot(ctx, "PAGE")

		books := make([]*Book, 0, len(bookNodes))

		for _, bookNode := range bookNodes {
			var ok bool
			var thumbnail, url, title, price, instockText string
			var ratingElemClass string

			timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			err := chromedp.Run(timeoutCtx,
				chromedp.AttributeValue(".image_container a img", "src", &thumbnail, &ok, chromedp.ByQuery, chromedp.FromNode(bookNode)),
				chromedp.AttributeValue(".image_container a", "href", &url, &ok, chromedp.ByQuery, chromedp.FromNode(bookNode)),
				chromedp.Text("h3 > a", &title, chromedp.ByQuery, chromedp.FromNode(bookNode)),
				chromedp.AttributeValue(".star-rating", "class", &ratingElemClass, &ok, chromedp.ByQuery, chromedp.FromNode(bookNode)),
				chromedp.Text("p.price_color", &price, chromedp.ByQuery, chromedp.FromNode(bookNode)),
				chromedp.Text("p.instock", &instockText, chromedp.ByQuery, chromedp.FromNode(bookNode)),
			)
			cancel()
			if errors.Is(context.DeadlineExceeded, err) {
				takeScreenShot(ctx, "timeout waiting for next button")
				log.Printf("Parse book timeout: %v\n", err)
				break
			} else if err != nil {
				log.Printf("Book Node: %+v\n", bookNode.Dump("", "  ", true))
				log.Fatal("Parse book item error: ", err)
			}

			books = append(books, &Book{
				Thumbnail: thumbnail,
				DetailURL: url,
				Title:     title,
				Rating:    getRating(ratingElemClass),
				Price:     price,
				Instock:   strings.Contains("In stock", strings.TrimSpace(instockText)),
			})
		}

		for _, book := range books {
			fmt.Printf("%v\n", book)
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		nextPageSelector := ".pager .next a"
		err = chromedp.Run(timeoutCtx,
			chromedp.WaitVisible(nextPageSelector, chromedp.ByQueryAll),
			chromedp.Click(nextPageSelector, chromedp.NodeVisible),
		)
		cancel()
		if errors.Is(context.DeadlineExceeded, err) {
			takeScreenShot(ctx, "timeout waiting for next button")
			break
		} else if err != nil {
			log.Fatal(`Wait for "next" button error: `, err)
		}
	}
}

type Book struct {
	Thumbnail string `json:"thumbnail"`
	DetailURL string `json:"detailURL"`
	Title     string `json:"title"`
	Rating    int    `json:"rating"`
	Price     string `json:"price"`
	Instock   bool   `json:"instock"`
}
