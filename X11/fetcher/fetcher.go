package fetcher

import (
	"bufio"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
)

func Fetch(link string) ([]byte, error) {
	random := browser.Random()
	client := resty.New()
	resp, err := client.R().
		SetHeader("User-Agent", random).
		Get(link)
	if err != nil {
		panic(err)
	}
	newReader := bufio.NewReader(resp.RawBody())
	e := determineEncoding(newReader)
	utf8Reader := transform.NewReader(newReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
