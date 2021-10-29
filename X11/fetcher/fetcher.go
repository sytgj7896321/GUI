package fetcher

import (
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

func Fetch(link string) ([]byte, error) {
	client := resty.New()
	client.RetryCount = 3
	client.RetryMaxWaitTime = 3 * time.Second
	resp, err := client.R().Get(link)
	if err != nil {
		log.Printf("Fetch %s Failed\n", link)
		return nil, err
	}
	return resp.Body(), nil
}
