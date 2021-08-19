package gcloud

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var userdir = "user/"
var projectdir = "project/"
var imagedir = "image/"
var videodir = "video/"
var betadir = "beta/"
var origindir = "origin/"

//GetBucket is function to get google cloud bucket.
func GetBucket(ctx context.Context) (*storage.BucketHandle, error) {
	config := &firebase.Config{
		ProjectID:     "sideproject-308610",
		StorageBucket: "sideproject-308610.appspot.com",
	}

	opt := option.WithCredentialsFile("../public/sideproject-308610-firebase-adminsdk-wmpyn-8b162af899.json")

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Println("[ERR] firebase err : ", err)
		return nil, err
	}

	client, err := app.Storage(ctx)
	if err != nil {
		log.Println("[ERR] storage err : ", err)
		return nil, err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Println("[ERR] default bucket err : ", err)
		return nil, err
	}

	return bucket, nil
}
