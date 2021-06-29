package common

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

type ProjectFile struct {
}

func (pf *ProjectFile) SaveBetaFile(ctx context.Context, userid *int, projectid int64, betafd *os.File) (string, error) {
	var id = uuid.New()

	bucket, err := GetBucket(ctx)
	if err != nil {
		log.Println("[ERR] failed to get bucket err : ", err)
		return "", err
	}

	storagedir := fmt.Sprintf("%s%d%s%s%d%s%s", userdir, userid, "/", projectdir, projectid, "/", betadir)

	object := bucket.Object(storagedir + betafd.Name())
	writer := object.NewWriter(ctx)

	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}

	_, err = io.Copy(writer, betafd)
	if err != nil {
		log.Println("[ERR] io copy err : ", err)
		return "", err
	}

	return storagedir + betafd.Name(), nil
}

//DeleteBetaFile is function to delete wrong betafile.
func DeleteBetaFile(ctx context.Context, userid *int, projectid int64, database *sql.DB) error {
	prefix := fmt.Sprintf("%s%d%s%s%d%s%s", userdir, userid, "/", projectdir, projectid, "/", betadir)
	var delimeter = "/"

	bucket, err := GetBucket(ctx)
	if err != nil {
		log.Println("[ERR] failed to get bucket err : ", err)
		return err
	}
	it := bucket.Objects(ctx, &storage.Query{Prefix: prefix, Delimiter: delimeter})

	betalink, err := ReadBetaLink(projectid, database)
	if err != nil {
		log.Println("[ERR] failed read beta link err : ", err)
		return err
	}

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

		if betalink != attrs.Name {
			err = bucket.Object(attrs.Name).Delete(ctx)
			if err != nil {
				log.Println("[ERR] failed to delete object err : ", err)
				return err
			}
		}
	}

	return nil
}

func ReadBetaLink(projectid int64, database *sql.DB) (string, error) {
	var betalink string

	stmt, err := database.Prepare(`SELECT
			beta_link 	
			FROM project
			WHERE id = ?`)
	if err != nil {
		log.Println("[ERR] prepare statement err : ", err)
		return "", err
	}
	defer stmt.Close()

	rows, err := stmt.Query(projectid)
	if err != nil {
		log.Println("[ERR] statement query err : ", err)
		return "", err
	}

	for rows.Next() {
		err = rows.Scan(&betalink)
		if err != nil {
			log.Println("[ERR] rows scan err : ", err)
			return "", err
		}
	}

	return betalink, nil
}

func (pf *ProjectFile) SaveOriginFile(ctx context.Context, userid *int, projectid int64, originfd *os.File) (string, error) {
	var id = uuid.New()
	bucket, err := GetBucket(ctx)
	if err != nil {
		log.Println("[ERR] failed to get bucket err : ", err)
		return "", err
	}

	storagedir := fmt.Sprintf("%s%d%s%s%d%s%s", userdir, userid, "/", projectdir, projectid, "/", origindir)

	object := bucket.Object(storagedir + originfd.Name())
	writer := object.NewWriter(ctx)

	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}

	_, err = io.Copy(writer, originfd)
	if err != nil {
		log.Println("[ERR] io copy err : ", err)
		return "", err
	}

	return storagedir + originfd.Name(), nil
}

//DeleteOriginFile is function to delete original file.
func DeleteOriginFile(ctx context.Context, userid *int, projectid int64, database *sql.DB) error {
	prefix := fmt.Sprintf("%s%d%s%s%d%s%s", userdir, userid, "/", projectdir, projectid, "/", origindir)
	var delimeter = "/"

	bucket, err := GetBucket(ctx)
	if err != nil {
		log.Println("[ERR] failed to get bucket err : ", err)
		return err
	}

	it := bucket.Objects(ctx, &storage.Query{Prefix: prefix, Delimiter: delimeter})

	originlink, err := ReadOriginLink(projectid, database)
	if err != nil {
		log.Println("[ERR] failed to read original file link err : ", err)
		return err
	}

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
		if originlink != attrs.Name {
			err = bucket.Object(attrs.Name).Delete(ctx)
			if err != nil {
				log.Println("[ERR] failed to delete object err : ", err)
				return err
			}
		}
	}

	return nil
}

//ReadOriginLink is function to read origin file link.
func ReadOriginLink(projectid int64, database *sql.DB) (string, error) {
	var originlink string

	stmt, err := database.Prepare(`SELECT
	origin_link
	FROM project WHERE id = ?`)
	defer stmt.Close()

	if err != nil {
		log.Println("[ERR] prepare statement err : ", err)
		return "", err
	}

	rows, err := stmt.Query(projectid)
	if err != nil {
		log.Println("[ERR] statement query err : ", err)
		return "", err
	}

	for rows.Next() {
		err = rows.Scan(&originlink)
		if err != nil {
			log.Println("[ERR] rows scan err : ", err)
			return "", err
		}
	}

	return originlink, nil
}
