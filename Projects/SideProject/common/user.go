package common

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
)

//ReadUserInfo is fuction to read user information.
func (m *mariadbHandler) ReadUserInfo(sessionid string) *ResUserInfo {
	var resuser *ResUserInfo

	resuser = &ResUserInfo{}

	stmt, err := m.db.Prepare(`SELECT id, nickname
				  FROM user WHERE session_id = ?`)
	if err != nil {
		log.Println("[ERR] prepare stmt err : ", err)
		return nil
	}

	defer stmt.Close()

	rows, err := stmt.Query(sessionid)
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

//UpdateUserInfo is function to update user information.
func (m *mariadbHandler) UpdateUserInfo(imglink string, reqjoininfo *ReqJoinInfo) error {

	tx, err := m.db.Begin()
	if err != nil {
		log.Println("[ERR] begin err : ", err)
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE user SET name=?, nickname=?, email=?,image_link=?,introduction=?,bank=?,account=?,updated_at=NOW()
	WHERE id=?`)
	if err != nil {
		log.Println("[ERR] transaction err : ", err)
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(reqjoininfo.UserInfo.Name, reqjoininfo.UserInfo.Nickname, reqjoininfo.UserInfo.Email, imglink, reqjoininfo.UserInfo.Introduction, reqjoininfo.AccountInfo.Bank, reqjoininfo.AccountInfo.Account, reqjoininfo.UserInfo.ID)
	if err != nil {
		log.Println("[ERR] Exec err:", err)
		return err
	}

	rowcnt, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERR] rows affected err:", err)
		return err
	}

	log.Println("[LOG] affected rows count :", rowcnt)

	err = tx.Commit()
	if err != nil {
		log.Println("[ERR] commit err : ", err)
		return err
	}

	return nil
}

//GetJoinInfo is function to get join information inserted.
func GetJoinInfo(multipartreader *multipart.Reader) (*ReqJoinInfo, []*os.File, error) {
	var reqjoininfo ReqJoinInfo
	var filearray []*os.File
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
		case "media":
			file, err := os.Create(part.Header.Get("Content-Filename"))
			if err != nil {
				log.Println("[ERR] file create err : ", err)
				return nil, nil, err
			}
			_, err = file.Write(data)
			if err != nil {
				log.Println("[ERR] file write err : ", err)
				return nil, nil, err
			}
			filearray = append(filearray, file)
		}
	}
	return &reqjoininfo, filearray, nil
}
