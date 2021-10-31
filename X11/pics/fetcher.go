package pics

import (
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

func Fetch(link string) ([]byte, error) {
	client := resty.New()
	client.RetryCount = 2
	client.RetryMaxWaitTime = 2 * time.Second
	resp, err := client.R().Get(link)
	if err != nil {
		log.Printf("Connect to source website %s failed, please check your network\n", link)
		return nil, err
	}
	return resp.Body(), nil
}
