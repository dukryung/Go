package common

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
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

//ReadProjectImageLinks is function to select project image links.
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

//SaveUserImgFile is fuction to save user image to google cloud stroage.
func SaveUserImgFile(args ArgsUpdateJoinUserInfo) (string, error) {
	id := uuid.New()
	bucket, err := GetBucket(args.ctx)
	if err != nil {
		log.Println("[ERR] failed to get bucket err : ", err)
		return "", err
	}

	storagedir := fmt.Sprintf("%s%d%s%s", userdir, args.joinuserinfo.UserInfo.UserID, "/", imagedir)

	object := bucket.Object(storagedir + args.userimgfile.Name())
	writer := object.NewWriter(args.ctx)

	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}

	_, err = io.Copy(writer, args.userimgfile)
	if err != nil {
		log.Println("[ERR] io copy err : ", err)
		return "", err
	}

	writer.Close()

	return storagedir + args.userimgfile.Name(), nil
}

//DeleteUserImgFile is function to check thage the user image folder exists in storage.
func DeleteUserImgFile(ctx context.Context, userid int64, database *sql.DB) error {

	var userimglink string
	stmt, err := database.Prepare(`SELECT image_link FROM user WHERE id = ?`)
	if err != nil {
		log.Println("[ERR] prepare statement err : ", err)
		return err
	}

	defer stmt.Close()

	rows, err := stmt.Query(userid)
	if err != nil {
		log.Println("[ERR] statement query err :", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&userimglink)
		if err != nil {
			log.Println("[ERR] rows scan err : ", err)
			return err
		}
	}

	bucket, err := GetBucket(ctx)
	if err != nil {
		log.Println("[ERR] failed to get bucket err : ", err)
		return err
	}

	var delimeter = "/"
	it := bucket.Objects(ctx, &storage.Query{Prefix: userimglink, Delimiter: delimeter})

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

		if userimglink != attrs.Name {
			err = bucket.Object(attrs.Name).Delete(ctx)
			if err != nil {
				log.Println("[ERR] failed to delete object err : ", err)
				return err
			}
		}
	}

	return nil
}
