package pics

import (
	"fmt"
	"github.com/cavaliercoder/grab"
	"os"
	"time"
)

type Task struct {
	ImageId *string
	Link    *string
	Percent float64
	Done    func()
	Cancel  func()
	Retry   func()
}

var (
	TaskList    = make(chan *Task, 1)
	In          = make(chan *grab.Request, 4)
	Out         = make(chan *grab.Response, 4)
	PercentList []*float64
)

//func AdvancedDownloadOriginal(task *Task) {
//	client := grab.NewClient()
//	req, _ := grab.NewRequest(LocalSaveDirectory, *task.Link)
//	resp := client.Do(req)
//	if err := resp.Err(); err != nil {
//		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
//		return
//	}
//	t := time.NewTicker(100 * time.Millisecond)
//	defer t.Stop()
//
//Loop:
//	for {
//		select {
//		case <-t.C:
//			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
//				resp.BytesComplete(),
//				resp.Size,
//				100*resp.Progress())
//		case <-resp.Done:
//			break Loop
//		}
//	}
//	//limiter := rate.NewLimiter(3, 5)
//	//req.RateLimiter = limiter
//	//PercentList = append(PercentList, &t.Percent)
//}

func DownloadOriginal() {
	for i := 0; i < 4; i++ {
		go makeProgressBar()
	}
	client := grab.NewClient()
	for i := 0; i < 4; i++ {
		go client.DoChannel(In, Out)
	}
}

func makeProgressBar() {
	for {
		select {
		case resp := <-Out:
			tick := time.NewTicker(25 * time.Millisecond)
		Loop:
			for {
				select {
				case <-tick.C:
					fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
						resp.BytesComplete(),
						resp.Size,
						100*resp.Progress())
				case <-resp.Done:
					fmt.Println("Done")
					tick.Stop()
					break Loop
				}
			}
			if err := resp.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
				os.Exit(1)
			}
		}
	}

}
