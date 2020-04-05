package jobutils

type Job struct {
	Post *Post  `json:"post"`
	Tag  string `json:"tag"`
	Done bool   `json:"done"`
}

func NewJob(number int, title, url, tag string, done bool) *Job {
	return &Job{
		Post: &Post{
			Number: number,
			Title:  title,
			URL:    url,
		},
		Tag:  tag,
		Done: done,
	}
}

type Post struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}

type Store interface {
	GetWaitingJobsByTag(string) ([]*Job, error)
	MakeJobDone(job *Job) error
	SubmitJobIfNotExist(job *Job) (bool, error)
	Close() error
}
