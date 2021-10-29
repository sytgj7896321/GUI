package pics

type Task struct {
	ImageId      int
	DownloadLink string
	TotalSize    int
	Cancel       func()
}
