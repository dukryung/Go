package common

import (
	"cloud.google.com/go/storage"
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

//Firebase object is to bind firebase-related functions.
type Firebase struct {
}

//CreateUserImages function is to insert images into firebase's storage.
func (f *Firebase) CreateUserImages(r *http.Request) error {

	config := &firebase.Config{
		ProjectID:     "sideproject-308610",
		StorageBucket: "sideproject-308610.appspot.com",
	}

	opt := option.WithCredentialsFile("./sideproject-308610-firebase-adminsdk-wmpyn-8b162af899.json")

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Println("[ERR] failed accessing firebase newapp : ", err)
		return err
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		log.Println("[ERR] faild firebase storage : ", err)
		return err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Println("[ERR] default bucket err : ", err)
		return err
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println("[ERR] parse MultipartForm err : ", err)
		return err
	}

	images := r.MultipartForm.File["user-images"]

	for _, image := range images {
		imagefd, err := image.Open()
		if err != nil {
			log.Println("[ERR] image files open err : ", err)
			return err
		}
		defer imagefd.Close()
	}

	bucket.Objects(context.Background(), storage.)

	

	

	return nil
}
