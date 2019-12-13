package jobutils

import (
	"context"
	"path"
	"strconv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Job struct {
	Post *Post  `json:"post"`
	Tag  string `json:"tag"`
	Done bool   `json:"done"`
}

type Post struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}

func NewFirestoreClient(projectID string) (*firestore.Client, error) {
	conf := &firebase.Config{
		ProjectID: projectID,
	}
	app, err := firebase.NewApp(context.Background(), conf)
	if err != nil {
		return nil, err
	}
	return app.Firestore(context.Background())
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

func GetWaitingJobsByTag(client *firestore.Client, tag string) ([]*Job, error) {
	jobs := []*Job{}

	iter := client.Collection(path.Join("jobs", "tags", tag)).Where("Done", "==", false).Documents(context.Background())

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		job := &Job{}
		if err := doc.DataTo(job); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func MakeJobDone(client *firestore.Client, job *Job) error {
	doneJob := job
	doneJob.Done = true
	_, err := client.Collection(path.Join("jobs", "tags", job.Tag)).Doc(strconv.Itoa(job.Post.Number)).Set(
		context.Background(),
		doneJob,
	)
	return err
}

func SubmitJobIfNotExist(client *firestore.Client, job *Job) error {
	snapShot, err := client.Collection(path.Join("jobs", "tags", job.Tag)).Doc(strconv.Itoa(job.Post.Number)).Get(context.Background())
	if err != nil && grpc.Code(err) != codes.NotFound {
		return err
	}
	if snapShot.Exists() {
		return nil
	}
	_, err = client.Doc((path.Join("jobs", "tags", job.Tag, strconv.Itoa(job.Post.Number)))).Set(
		context.Background(),
		job,
	)
	return err
}
