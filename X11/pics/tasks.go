package pics

type Task struct {
	ImageId      int
	DownloadLink string
	TotalSize    int
	Cancel       func()
	Retry        func()
}

var (
	taskList = make(chan *Task, 48)
)
