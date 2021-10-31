package pics

import "github.com/cavaliercoder/grab"

type Task struct {
	ImageId     *string
	Link        *string
	CurrentSize int
	TotalSize   int
	Cancel      func()
	Retry       func()
}

var (
	TaskList []*Task
)

func (t *Task) GetOriginal(dst string) {
	client := grab.NewClient()
	req, _ := grab.NewRequest(dst, *t.Link)
	client.Do(req)
}
