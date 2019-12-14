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

type fireStore struct {
	client *firestore.Client
}

func NewFireStore(projectID string) (Store, error) {
	conf := &firebase.Config{
		ProjectID: projectID,
	}
	app, err := firebase.NewApp(context.Background(), conf)
	if err != nil {
		return nil, err
	}
	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, err
	}
	return &fireStore{client: client}, nil
}

func (f *fireStore) GetWaitingJobsByTag(tag string) ([]*Job, error) {
	jobs := []*Job{}

	iter := f.client.Collection(path.Join("jobs", "tags", tag)).Where("Done", "==", false).Documents(context.Background())

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

func (f *fireStore) MakeJobDone(job *Job) error {
	doneJob := job
	doneJob.Done = true
	_, err := f.client.Collection(path.Join("jobs", "tags", job.Tag)).Doc(strconv.Itoa(job.Post.Number)).Set(
		context.Background(),
		doneJob,
	)
	return err
}

func (f *fireStore) SubmitJobIfNotExist(job *Job) error {
	snapShot, err := f.client.Collection(path.Join("jobs", "tags", job.Tag)).Doc(strconv.Itoa(job.Post.Number)).Get(context.Background())
	if err != nil && grpc.Code(err) != codes.NotFound {
		return err
	}
	if snapShot.Exists() {
		return nil
	}
	_, err = f.client.Doc((path.Join("jobs", "tags", job.Tag, strconv.Itoa(job.Post.Number)))).Set(
		context.Background(),
		job,
	)
	return err
}

func (f *fireStore) Close() error {
	return f.client.Close()
}
