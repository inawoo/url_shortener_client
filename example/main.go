package main

import (
	"fmt"
	"github.com/inawoo/url_shortener/client"
)

func main() {

	urlShortener := client.NewClient(client.WithBaseURL("https://staging.link.inawo.live"), client.WithPoolCount(10))

	_, er := urlShortener.CheckHealth()
	if er != nil {
		panic(er)
	}

	for i := 0; i < 100; i++ {

		response, err := urlShortener.ShortenURL(client.ShortenURLRequest{
			URL: "https://www.google.com/search?q=hello+world&oq=hello+world&aqs=chrome..69i57j0l7.1001j0j7&sourceid=chrome&ie=UTF-" + fmt.Sprint(i),
		})
		if err != nil {
			panic(err)
		}

		fmt.Println(response)
	}

}
