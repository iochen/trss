package trss

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

// GenerateRSS generates RSS feed from given channel name
// NOTE: arg [name] must NOT contain an "@"
func GenerateRSS(name string) (feeds.Rss, error) {
	link := fmt.Sprintf("https://t.me/s/%s", name)
	rss := feeds.Rss{
		Feed: &feeds.Feed{
			Link: &feeds.Link{
				Href: link,
			},
		},
	}

	resp, err := http.Get(link)
	if err != nil {
		return feeds.Rss{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return feeds.Rss{}, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return feeds.Rss{}, err
	}

	title, err := document.Find(".tgme_header_title").Html()
	if err != nil {
		return feeds.Rss{}, err
	}

	rss.Feed.Title = title

	document.Find(".js-widget_message").
		Each(func(i int, selection *goquery.Selection) {
			item := &feeds.Item{}
			dateSec := selection.Find(".tgme_widget_message_date")
			link, _ := dateSec.Attr("href")
			item.Link = &feeds.Link{
				Href: link,
			}

			datetime, _ := dateSec.Find("time").Attr("datetime")
			item.Created, err = time.Parse(time.RFC3339, datetime)
			if err != nil {
				return
			}

			text, err := selection.Find(".js-message_text").Html()
			if err != nil {
				return
			}

			item.Content = text

			rss.Add(item)
		})

	return rss, nil
}
