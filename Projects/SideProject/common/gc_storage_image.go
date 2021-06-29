package common

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/api/iterator"
	"google.golang.org/appengine"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//Image is struct to make  image function concerned  objectively.
type Image struct {
}

//SaveProjectImages is function to save project images.
func (i *Image) SaveProjectImages(ctx context.Context, userid *int, projectid int64, imagefdarr []*os.File) ([]string, error) {
	var id = uuid.New()
	var imagefullpatharr []string

	bucket, err := GetBucket(ctx)
	if err != nil {
		log.Println("[ERR] failed to get bucket err : ", err)
		return nil, err
	}

	storagedir := fmt.Sprintf("%s%d%s%s%d%s", userdir, userid, "/", projectdir, projectid, imagedir)
	for _, imagefd := range imagefdarr {
		object := bucket.Object(storagedir + imagefd.Name())
		writer := object.NewWriter(ctx)
		writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}
		_, err = io.Copy(writer, imagefd)
		if err != nil {
			log.Println("[ERR] io copy err : ", err)
			return nil, err
		}
		imagefullpatharr = append(imagefullpatharr, storagedir+imagefd.Name())
	}

	return imagefullpatharr, nil
}

//DeleteProjectImages is function to delete project images.
func DeleteProjectImages(ctx context.Context, userid *int, projectid int64, database *sql.DB) error {
	var prefix = fmt.Sprintf("%s%d%s%s%d%s", userdir, userid, "/", projectdir, projectid, imagedir)
	var delimeter = "/"

	bucket, err := GetBucket(ctx)
	if err != nil {
		log.Println("[ERR] failed to get bucket err : ", err)
		return err
	}

	it := bucket.Objects(ctx, &storage.Query{Prefix: prefix, Delimiter: delimeter})

	projectimagelinks, err := ReadProjectImageLinks(projectid, database)
	if err != nil {
		log.Println("[ERR] failed to select project image links err : ", err)
		return err
	}

	for {
		var isprojectimages = false
		attrs, err := it.Next()
		if err == iterator.Done {
			log.Println("[LOG] iterator done!")
			break
		}

		if err != nil {
			log.Println("[ERR] bucket objects err : ", err)
			return err
		}

		for _, projectimagelink := range projectimagelinks {
			if attrs.Name == projectimagelink {
				isprojectimages = true
			}
		}

		if isprojectimages == false {
			err = bucket.Object(attrs.Name).Delete(ctx)
			if err != nil {
				log.Println("[ERR] failed to delete object err : ", err)
				return err
			}
		}
	}

	return nil
}

//SelectProjectImageLinks is function to select project image links.
func ReadProjectImageLinks(projectid int64, database *sql.DB) ([]string, error) {
	var projectimagelinks []string
	var projectimagelink string

	stmt, err := database.Prepare(`SELECT
					link
					FROM image 
					WHERE project_id = ?`)
	if err != nil {
		log.Println("[ERR] prepare satement err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(projectid)
	if err != nil {
		log.Println("[ERR] query err : ", err)
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&projectimagelink)
		if err != nil {
			log.Println("[ERR] rows scan err : ", err)
			return nil, err
		}

		projectimagelinks = append(projectimagelinks, projectimagelink)
	}

	return projectimagelinks, nil
}

//SaveImageFiles is function to get user's images
func (i *Image) SaveImageFiles(c *gin.Context) (string, *ReqJoinInfo, error) {

	multipartreader, err := c.Request.MultipartReader()
	if err != nil {
		log.Println("[ERR] multipartreader err : ", err)
		return "", nil, err
	}

	reqjoininfo, filearray, err := GetJoinInfo(multipartreader)
	if err != nil {
		log.Println("[ERR] getjoininfo err : ", err)
		return "", nil, err
	}

	ctx := appengine.NewContext(c.Request)
	var filefullpath string
	for _, file := range filearray {

		filefullpath, err = CreateProfileImage(ctx, file, file.Name(), reqjoininfo.UserInfo.ID)
		if err != nil {
			log.Println("[ERR] create image file err : ", err)
			err := DeleteProfileImage(ctx, reqjoininfo.UserInfo.ID)
			if err != nil {
				log.Println("[ERR] failed delete images err : ", err)
				return "", nil, err
			}
			return "", nil, err
		}
		file.Close()
	}
	return filefullpath, reqjoininfo, nil
}

//CreateProfileImage is Function to post a imagefile for Google Storage cloud.
func CreateProfileImage(ctx context.Context, file *os.File, filename string, userid string) (string, error) {

	id := uuid.New()
	bucket, err := GetBucket(ctx)
	if err != nil {
		log.Println("[ERR] failed to get bucket err : ", err)
		return "", err
	}

	storagedir := userdir + userid + "/" + imagedir

	object := bucket.Object(storagedir + filename)
	writer := object.NewWriter(ctx)

	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}

	_, err = io.Copy(writer, file)
	if err != nil {
		log.Println("[ERR] io copy err : ", err)
		return "", err
	}

	writer.Close()

	return storagedir + filename, nil
}

//DeleteProfileImage is function to check thage the user image folder exists in storage.
func DeleteProfileImage(ctx context.Context, userid string) error {
	var prefix = userdir + userid + "/" + imagedir
	var delimeter = "/"

	bucket, err := GetBucket(ctx)
	if err != nil {
		log.Println("[ERR] failed to get bucket err : ", err)
		return err
	}

	it := bucket.Objects(ctx, &storage.Query{Prefix: prefix, Delimiter: delimeter})

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

		err = bucket.Object(attrs.Name).Delete(ctx)
		if err != nil {
			log.Println("[ERR] failed to delete object err : ", err)
			return err
		}
	}

	return nil
}
