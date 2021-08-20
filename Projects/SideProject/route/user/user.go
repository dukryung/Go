package user

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sideproject/route/common"
	"sideproject/route/gcloud"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

type Usr struct {
	DB *sql.DB
}

// JSUser is used to indicate that these are json type's artist informations.
type JSUser struct {
	UserID              int64  `json:"user_id" binding:"required"`
	Name                string `json:"name"`
	Nickname            string `json:"nickname"`
	Email               string `json:"email"`
	Introduction        string `json:"introduction"`
	AgreeEmailMarketing bool   `json:"agree_email_marketing"`
}

//ResUserInfo is structure to contain user information in index page.
type ResJSUser struct {
	ID       string `json:"id"`
	NickName string `json:"nickname"`
}

//ArgsUpdateJoinUserInfo is a collection of extracted prameters that is user's information.

type ReqJoinInfo struct {
	UserInfo    JSUser  `json:"user"`
	AccountInfo Account `json:"account"`
}

type Account struct {
	UserID      int64  `json:"user_id"`
	Bank        int    `json:"bank"`
	Account     string `json:"account"`
	AgreePolicy bool   `json:"agree_policy"`
}

type AuthUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Social        string `json:"social"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func (u *Usr) Routes(route *gin.RouterGroup) {
	route.GET("/", u.getUser)
	route.PUT("/", u.putUser)
}

func (u *Usr) getUser(c *gin.Context) {

	jsu := &JSUser{}

	err := c.ShouldBindJSON(jsu)
	if err != nil {
		log.Println("[ERR] failed to extract user idd err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	resuser := u.ReadUserInfo(jsu.UserID)
	c.JSON(http.StatusOK, resuser)
}

func (u *Usr) putUser(c *gin.Context) {
	filefullpath, err := u.SaveJoinUserInfo(c)
	if err != nil {
		log.Println("[ERR] get image file err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filedirinfo": filefullpath,
	})
}

//ReadUserInfo is fuction to read user information.
func (u *Usr) ReadUserInfo(userid int64) *ResJSUser {

	var resjsu = &ResJSUser{}

	stmt, err := u.DB.Prepare(`SELECT id, nickname
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
		err = rows.Scan(&resjsu.ID, &resjsu.NickName)
		if err != nil {
			log.Println("[ERR] stmt query err : ", err)
			return nil
		}
	}

	return resjsu
}

func (u *Usr) ReadUserID(sessionid string) (*int, error) {

	var userid *int

	stmt, err := u.DB.Prepare(`SELECT id FROM user WHERE session_id = ?`)
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
func (u *Usr) CreateUserInfo(authuserinfo AuthUserInfo) error {
	tx, err := u.DB.Begin()
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
func (u *Usr) SaveJoinUserInfo(c *gin.Context) (string, error) {

	var args common.ArgsUpdateJoinUserInfo
	multipartreader, err := c.Request.MultipartReader()
	if err != nil {
		log.Println("[ERR] multipartreader err : ", err)
		return "", err
	}

	args.Joinuserinfo, args.Userimgfile, err = ExtractJoinUserInfo(multipartreader)
	if err != nil {
		log.Println("[ERR] getjoininfo err : ", err)
		return "", err
	}

	defer args.Userimgfile.Close()

	args.CTX = appengine.NewContext(c.Request)
	args.DB = u.DB

	filefullpath, err := UpdateJoinUserInfo(args)
	if err != nil {
		log.Println("[ERR] create image file err : ", err)

		return "", err
	}

	return filefullpath, nil
}

//UpdateJoinUserInfo is Function to post a imagefile for Google Storage cloud.
func UpdateJoinUserInfo(args common.ArgsUpdateJoinUserInfo) (string, error) {

	tx, err := args.DB.Begin()
	if err != nil {
		log.Println("[ERR] begin err : ", err)
		return "", err
	}

	defer gcloud.DeleteUserImgFile(args.CTX, args.Joinuserinfo.UserInfo.UserID, args.DB)
	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE user SET name=?, nickname=?, email=?,image_link=?,introduction=?,bank=?,account=?,updated_at=NOW()
	WHERE id=?`)
	if err != nil {
		log.Println("[ERR] transaction err : ", err)
		return "", err
	}
	defer stmt.Close()

	savedimagepath, err := gcloud.SaveUserImgFile(args)
	if err != nil {
		log.Println("[ERR] failed to save user image file to google cloud storage err : ", err)
		return "", err
	}

	if err != nil {
		log.Println("[ERR] Exec err:", err)
		return "", err
	}
	result, err := stmt.Exec(args.Joinuserinfo.UserInfo.Name, args.Joinuserinfo.UserInfo.Nickname, args.Joinuserinfo.UserInfo.Email, savedimagepath, args.Joinuserinfo.UserInfo.Introduction, args.Joinuserinfo.AccountInfo.Bank, args.Joinuserinfo.AccountInfo.Account, args.Joinuserinfo.UserInfo.UserID)

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
func ExtractJoinUserInfo(multipartreader *multipart.Reader) (*common.ReqJoinInfo, *os.File, error) {
	var reqjoininfo common.ReqJoinInfo
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
