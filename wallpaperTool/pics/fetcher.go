package pics

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/go-resty/resty/v2"
	"time"
)

func Fetch(link string) ([]byte, error) {
	client := resty.New()
	client.RetryCount = 3
	client.RetryMaxWaitTime = 3 * time.Second
	resp, err := client.R().Get(link)
	if err != nil {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Wallpaper Tool",
			Content: fmt.Sprintf("Connect to source website %s failed, please check your network\n", link),
		})
		return nil, err
	}
	return resp.Body(), nil
}
