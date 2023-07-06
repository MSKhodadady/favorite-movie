package myPac

import (
	"net/url"

	"github.com/gocolly/colly/v2"
)

//: section[data-testid="find-results-section-title"]>div>ul>li
//: div>div>a
//: div>div>ul>li>span

func SearchForMovieName(txt string) []*Movie {

	movies := make([]*Movie, 0)

	collector := colly.NewCollector()

	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US")
	})

	collector.OnHTML("section[data-testid=\"find-results-section-title\"]"+
		">div>ul>li", func(h *colly.HTMLElement) {

		f := &Movie{}
		err := h.Unmarshal(f)

		if err != nil {
			panic(err)
		}

		movies = append(movies, f)
	})

	visitUrl := "https://www.imdb.com/find/?q=" + url.QueryEscape(txt)

	collector.Visit(visitUrl)

	return movies
}
