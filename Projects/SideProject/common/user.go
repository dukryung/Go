package common

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

type User struct {
	UserID              int64  `json:"user_id" binding:"required"`
	Name                string `json:"name"`
	Nickname            string `json:"nickname"`
	Email               string `json:"email"`
	Introduction        string `json:"introduction"`
	AgreeEmailMarketing bool   `json:"agree_email_marketing"`
}

type Account struct {
	UserID      int64  `json:"user_id"`
	Bank        int    `json:"bank"`
	Account     string `json:"account"`
	AgreePolicy bool   `json:"agree_policy"`
}

//ArgsUpdateJoinUserInfo is a collection of extracted prameters that is user's information.
type ArgsUpdateJoinUserInfo struct {
	ctx          context.Context
	userimgfile  *os.File
	joinuserinfo *ReqJoinInfo
	database     *sql.DB
}

//ReadUserInfo is fuction to read user information.
func (m *mariadbHandler) ReadUserInfo(userid int64) *ResUserInfo {
	var resuser *ResUserInfo

	resuser = &ResUserInfo{}

	stmt, err := m.db.Prepare(`SELECT id, nickname
				  FROM user WHERE user.id = ?`)
	if err != nil {
		log.Println("[ERR] prepare stmt err : ", err)
		return nil
	}

	defer stmt.Close()

	rows, err := stmt.Query(userid)
	if err != nil {
		log.Println("[ERR] stmt query err : ", err)
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&resuser.ID, &resuser.NickName)
		if err != nil {
			log.Println("[ERR] stmt query err : ", err)
			return nil
		}
	}

	return resuser
}

func (m *mariadbHandler) ReadUserID(sessionid string) (*int, error) {

	var userid *int

	stmt, err := m.db.Prepare(`SELECT id FROM user WHERE session_id = ?`)
	if err != nil {
		log.Println("[ERR] prepare stmt err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(sessionid)
	if err != nil {
		log.Println("[ERR] stmt query err : ", err)
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(userid)
		if err != nil {
			return nil, err
		}
	}

	return userid, nil

}

//CreateUserInfo is function to create user information.
func (m *mariadbHandler) CreateUserInfo(authuserinfo AuthUserInfo) error {
	tx, err := m.db.Begin()
	if err != nil {
		log.Println("[ERR] transaction begin err : ", err)
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT IGNORE INTO user (session_id,email,social) VALUES(?,?,?)`)
	if err != nil {
		log.Println("[ERR]  prepared statement err : ", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(authuserinfo.ID, authuserinfo.Email, authuserinfo.Social)
	if err != nil {
		log.Println("[ERR] exec statement err : ", err)
		return err
	}

	rowcnt, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERR] rows affected err : ", err)
		return err
	}
	if rowcnt == 0 {
		log.Println("[LOG] None affected rows")
	}

	err = tx.Commit()
	if err != nil {
		log.Println("[ERR] transacntion commit err : ", err)
		return err
	}

	return nil
}

//SaveImageFiles is function to get user's images
func (m *mariadbHandler) SaveJoinUserInfo(c *gin.Context) (string, error) {
	var args ArgsUpdateJoinUserInfo
	multipartreader, err := c.Request.MultipartReader()
	if err != nil {
		log.Println("[ERR] multipartreader err : ", err)
		return "", err
	}

	args.joinuserinfo, args.userimgfile, err = ExtractJoinUserInfo(multipartreader)
	if err != nil {
		log.Println("[ERR] getjoininfo err : ", err)
		return "", err
	}

	defer args.userimgfile.Close()

	args.ctx = appengine.NewContext(c.Request)
	args.database = m.db

	filefullpath, err := UpdateJoinUserInfo(args)
	if err != nil {
		log.Println("[ERR] create image file err : ", err)

		return "", err
	}

	return filefullpath, nil
}

//UpdateJoinUserInfo is Function to post a imagefile for Google Storage cloud.
func UpdateJoinUserInfo(args ArgsUpdateJoinUserInfo) (string, error) {

	tx, err := args.database.Begin()
	if err != nil {
		log.Println("[ERR] begin err : ", err)
		return "", err
	}

	defer DeleteUserImgFile(args.ctx, args.joinuserinfo.UserInfo.UserID, args.database)
	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE user SET name=?, nickname=?, email=?,image_link=?,introduction=?,bank=?,account=?,updated_at=NOW()
	WHERE id=?`)
	if err != nil {
		log.Println("[ERR] transaction err : ", err)
		return "", err
	}
	defer stmt.Close()

	savedimagepath, err := SaveUserImgFile(args)
	if err != nil {
		log.Println("[ERR] failed to save user image file to google cloud storage err : ", err)
		return "", err
	}

	result, err := stmt.Exec(args.joinuserinfo.UserInfo.Name, args.joinuserinfo.UserInfo.Nickname, args.joinuserinfo.UserInfo.Email, savedimagepath, args.joinuserinfo.UserInfo.Introduction, args.joinuserinfo.AccountInfo.Bank, args.joinuserinfo.AccountInfo.Account, args.joinuserinfo.UserInfo.UserID)
	if err != nil {
		log.Println("[ERR] Exec err:", err)
		return "", err
	}

	rowcnt, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERR] rows affected err:", err)
		return "", err
	}

	log.Println("[LOG] affected rows count :", rowcnt)

	err = tx.Commit()
	if err != nil {
		log.Println("[ERR] commit err : ", err)
		return "", err
	}

	return savedimagepath, nil
}

//ExtractJoinUserInfo is function to get join information inserted.
func ExtractJoinUserInfo(multipartreader *multipart.Reader) (*ReqJoinInfo, *os.File, error) {
	var reqjoininfo ReqJoinInfo
	var userimgfile *os.File

	for {
		part, err := multipartreader.NextPart()

		if err == io.EOF {
			log.Println("[LOG] part done")
			break
		}

		if err != nil {
			log.Println("[ERR] part err : ", err)
			return nil, nil, err
		}

		data, err := ioutil.ReadAll(part)
		if err != nil {
			log.Println("[ERR] ioutil read all err : ", err)
			return nil, nil, err
		}

		switch part.Header.Get("Content-ID") {
		case "metadata":
			err = json.Unmarshal(data, &reqjoininfo)
			if err != nil {
				log.Println("[ERR] hson unmarshal err : ", err)
				return nil, nil, err
			}
		case "userimage":
			userimgfile, err = os.Create(part.Header.Get("Content-Filename"))
			if err != nil {
				log.Println("[ERR] file create err : ", err)
				return nil, nil, err
			}
			_, err = userimgfile.Write(data)
			if err != nil {
				log.Println("[ERR] file write err : ", err)
				return nil, nil, err
			}

		}
	}
	return &reqjoininfo, userimgfile, nil
}
