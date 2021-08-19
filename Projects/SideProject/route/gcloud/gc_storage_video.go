package gcloud

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

type Video struct {
}

func (v *Video) SaveVideoFile(ctx context.Context, userid int64, projectid int64, videofd *os.File) (string, error) {
	var id = uuid.New()

	bucket, err := GetBucket(ctx)
	if err != nil {
		log.Println("[ERR] failed to get bucket err : ", err)
		return "", err
	}

	storagedir := fmt.Sprintf("%s%d%s%s%d%s%s", userdir, userid, "/", projectdir, projectid, "/", videodir)

	object := bucket.Object(storagedir + videofd.Name())
	writer := object.NewWriter(ctx)

	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}

	_, err = io.Copy(writer, videofd)
	if err != nil {
		log.Println("[ERR] io copy err : ", err)
		return "", err
	}

	return storagedir + videofd.Name(), nil
}

//DeleteVideoFile is function to delete video file.
func DeleteVideoFile(ctx context.Context, userid int64, projectid int64, database *sql.DB) error {
	prefix := fmt.Sprintf("%s%d%s%s%d%s%s", userdir, userid, "/", projectdir, projectid, "/", videodir)
	var delimeter = "/"

	bucket, err := GetBucket(ctx)
	if err != nil {
		log.Println("[ERR] failed to get bucket err : ", err)
		return err
	}

	it := bucket.Objects(ctx, &storage.Query{Prefix: prefix, Delimiter: delimeter})
	videolink, err := ReadVideoFileLink(projectid, database)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			log.Println("[LOG] iterator done!")
			break
		}

		if err != nil {
			log.Println("[ERR] bucket objects err : ", err)
			return err
		}
		if videolink != attrs.Name {
			err = bucket.Object(attrs.Name).Delete(ctx)
			if err != nil {
				log.Println("[ERR] failed to delete object err : ", err)
				return err
			}
		}
	}

	return nil
}

//ReadVideoFileLink is function to read video file link.
func ReadVideoFileLink(projectid int64, database *sql.DB) (string, error) {
	var videolink string
	stmt, err := database.Prepare(`SELECT 
					  video_link
					  FROM project 
					  WHERE id = ?`)

	defer stmt.Close()

	rows, err := stmt.Query(projectid)
	if err != nil {
		log.Println("[ERR] prepare statment query err : ", err)
		return "", err
	}

	for rows.Next() {
		err = rows.Scan(&videolink)
		if err != nil {
			log.Println("[ERR] rows scan err : ", err)
			return "", err
		}
	}

	return videolink, nil
}
