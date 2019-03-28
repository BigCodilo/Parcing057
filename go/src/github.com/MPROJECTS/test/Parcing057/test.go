package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

//News - struct of the news
type News struct {
	Title string
	Text  string
}

func main() {
	flagLoading := true
	go DotCounter(1, &flagLoading)
	url := "https://www.057.ua/news"
	document := GetDocument(url)
	news := GetNews(document)
	flagLoading = false
	for index, oneOfNews := range news {
		fmt.Println(index, ":", oneOfNews.Title, "\n\n", oneOfNews.Text, "\n\n\n\n_________________________________________________________")
	}
}

//GetDocument - return a document of HTML page (russian symbols)
func GetDocument(url string) *goquery.Document {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error")
	}

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return document
}

//GetNews - return a slice of news (title + text)
func GetNews(document *goquery.Document) []News {
	news := []News{}
	document.Find(".col-lg-9").Each(func(_ int, panelOfNewsBlock *goquery.Selection) {
		panelOfNewsBlock.Find(".c-news-card").Each(func(_ int, blockOfNews *goquery.Selection) {
			title := blockOfNews.Find(".c-news-card__title").Text()
			title = strings.Replace(title, ", - ФОТО", "", -1)
			title = strings.Replace(title, ", ФОТО", "", -1)
			title = strings.Replace(title, ", - ВИДЕО", "", -1)
			title = strings.Replace(title, ", - ВИДЕО, ФОТО", "", -1)

			url, _ := blockOfNews.Find("a").Attr("href")
			url = "https://www.057.ua" + url

			url = strings.Replace(url, "#comments", "", -1)
			dicumentNews := GetDocument(url)
			text := GetNewsBody(dicumentNews)
			news = append(news, News{
				title,
				text,
			})

		})
	})
	return news
}

//GetNewsBody - return a text from news
func GetNewsBody(document *goquery.Document) string {
	textAll := ""
	document.Find(".col-lg-9").Each(func(_ int, bodyOfNews *goquery.Selection) {
		bodyOfNews.Find(".article-text").Each(func(_ int, bodyOfBody *goquery.Selection) {
			text := bodyOfBody.Find("p").Text()
			textAll += text
		})
	})
	return textAll
}

//DotCounter - print a time of wait
func DotCounter(count int, flag *bool) {
	if !*flag {
		return
	}

	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Print("\n\n\n    ____parcing_of_data_____", count, "___second(s)____")

	count++
	time.Sleep(1 * time.Second)
	DotCounter(count, flag)
}
