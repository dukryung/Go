package common

import (
	"log"
)

type ResProfileFrameInfo struct {
	UserID           int64  `json:"user_id"`
	UserNickName     string `json:"user_nickname"`
	UserIntroduction string `json:"user_introduction"`
	CreatedAt        string `json:"created_at"`
	ProjectCount     int    `json:"project_cnt"`
	SellCount        int    `json:"sell_cnt"`
	BuyCount         int    `json:"buy_cnt"`
	WithdrawCount    int    `json:"withdraw_cnt"`
}

type ResProfileProjectInfo struct {
	ProjectList []Project `json:"project_list"`
}

type Project struct {
	ProjectID    int64  `json:"project_id"`
	Title        string `json:"title"`
	Desc         string `json:"desc"`
	ImageLink    string `json:"image_link"`
	CreatedAt    string `json:"created+at"`
	SellCount    int    `json:"sell_cnt"`
	UserNickName string `json:"user_nickname"`
	CommontCount int    `json:"comment_count"`
	UpvoteCount  int    `json:"upvote_count"`
	Price        int    `json:"price"`
	Beta         bool   `json:"beta"`
}

type ResProfileSellInfo struct {
	SellList []Sell `json:"sell_list"`
}

type Sell struct {
	SellID        int64  `json:"sell_id"`
	Title         string `json:"title"`
	Date          string `json:"date"`
	BuyerID       int    `json:"buyer_id"`
	BuyerNickName string `json:"buyer_nickname"`
	Price         int    `json:"price"`
}

type ResProfileBuyInfo struct {
	BuyList []Buy `json:"buy_list"`
}

type Buy struct {
	ID             int64  `json:"buy_id"`
	Title          string `json:"title"`
	Date           string `json:"date"`
	SellerID       int    `json:"selller_id"`
	SellerNickName string `json:"seller_nickname"`
	Price          int    `json:"price"`
}

type ResProfileWithdrawInfo struct {
	WithdrawList []Withdraw `json:"withdraw_list"`
}
type Withdraw struct {
	ID            int64  `json:"withdraw_id"`
	RequestedDate string `json:"requested_date"`
	CompleteDate  string `json:"complete_date"`
	Amount        int    `json:"amount"`
}

type ReqModificationUserInfo struct {
	UserID              int64  `json:"user_id"`
	Name                string `json:"name"`
	Nickname            string `json:"nickname"`
	Email               string `json:"email"`
	AgreeEmailMarketing bool   `json:"agree_email_marketing"`
	Introduction        string `json:"introduction"`
	ImageLink           string `json:"image_link"`
}

type ResModificationUserInfo struct {
	UserID              int64  `json:"user_id"`
	Name                string `json:"name"`
	Nickname            string `json:"nickname"`
	Email               string `json:"email"`
	AgreeEmailMarketing bool   `json:"agree_email_marketing"`
	Introduction        string `json:"introduction"`
	ImageLink           string `json:"image_link"`
}

type ResArtistProfileInfo struct {
	ArtistInfo  ArtistProfile `json:"artist"`
	ProjectList []Project     `json:"project_list"`
}

type ArtistProfile struct {
	ArtistID       int64  `json:"artist_id"`
	ArtistNickName string `json:"artist_nickname"`
	Introduction   string `json:"introduction"`
	ProjectCount   int    `json:"project_cnt"`
	SellCount      int    `json:"sell_cnt"`
	ImageLink      string `json:"image_link"`
}

type ResPersonalInformation struct {
	Account Account `json:"account"`
}
type ReqPersonalInformation struct {
	Account Account `json:"account"`
}

func (m *mariadbHandler) ReadProfileFrameInfo(userid int64) (*ResProfileFrameInfo, error) {

	var resprofileframeinfo = &ResProfileFrameInfo{}

	stmt, err := m.db.Prepare(`SELECT 
	u.id, 
	u.nickname,
	u.introduction,
	u.created_at,
	(SELECT COUNT(id) FROM project WHERE user_id = ?),
	(SELECT COUNT(id) FROM sell_history WHERE user_id = ?),
	(SELECT COUNT(id) FROM buy_history WHERE user_id = ?),
	(SELECT COUNT(id) FROM withdraw_history WHERE user_id = ?)
	FROM user AS u 
	WHERE id = ?;`)
	if err != nil {
		log.Println("[ERR] prepare stmt err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(userid, userid, userid, userid, userid)
	if err != nil {
		log.Println("[ERR] stmt query err : ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&resprofileframeinfo.UserID, &resprofileframeinfo.UserNickName, &resprofileframeinfo.UserIntroduction, &resprofileframeinfo.CreatedAt, &resprofileframeinfo.ProjectCount, &resprofileframeinfo.SellCount, &resprofileframeinfo.BuyCount, &resprofileframeinfo.WithdrawCount)
		if err != nil {
			log.Println("[ERR] stmt query err : ", err)
			return nil, err
		}
	}

	return resprofileframeinfo, nil
}

func (m *mariadbHandler) ReadProfileProjectInfo(userid int64) (*ResProfileProjectInfo, error) {
	var resprofileprojectinfo = &ResProfileProjectInfo{}
	var project Project

	stmt, err := m.db.Prepare(`SELECT 
	p.id, 
	p.title,
	p.description,
	i.link,
	p.created_at,
	(SELECT COUNT(id) FROM sell_history WHERE user_id = ?),
	p.comment_count,
	p.total_upvote_count,
	p.price,
	p.beta
	FROM user AS u 
	INNER JOIN project AS p ON p.user_id = u.id 
	INNER JOIN image AS i ON i.project_id = p.id 
	INNER JOIN 
	(SELECT project_id, MIN(created_at) created_at 
	FROM image GROUP BY project_id) AS ii ON ii.project_id = i.project_id AND i.created_at = ii.created_at 
	WHERE u.id = ? GROUP BY p.id;`)
	if err != nil {
		log.Println("[ERR] prepare stmt err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(userid, userid)
	if err != nil {
		log.Println("[ERR] stmt query err : ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&project.ProjectID, &project.Title, &project.Desc, &project.ImageLink, &project.CreatedAt, &project.SellCount, &project.CommontCount, &project.CommontCount, &project.Price, &project.Beta)
		if err != nil {
			log.Println("[ERR] stmt query err : ", err)
			return nil, err
		}

		resprofileprojectinfo.ProjectList = append(resprofileprojectinfo.ProjectList, project)
	}

	return resprofileprojectinfo, nil

}

func (m *mariadbHandler) ReadProfileSellInfo(userid int64) (*ResProfileSellInfo, error) {
	var resprofilesellinfo = &ResProfileSellInfo{}
	var sell Sell

	stmt, err := m.db.Prepare(`  SELECT 
							p.id, 
							p.title,
							sh.created_at,
							sh.buyer_id,
							sh.buyer_nickname,
							sh.price
							FROM user AS u
							INNER JOIN sell_history as sh ON sh.user_id = u.id 
							INNER JOIN project AS p ON p.id = sh.project_id 
							WHERE u.id= ?;
	`)
	if err != nil {
		log.Println("[ERR] prepare stmt err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(userid)
	if err != nil {
		log.Println("[ERR] stmt query err : ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&sell.SellID, &sell.Title, &sell.Date, &sell.BuyerID, &sell.BuyerNickName, &sell.Price)
		if err != nil {
			log.Println("[ERR] stmt query err : ", err)
			return nil, err
		}
		resprofilesellinfo.SellList = append(resprofilesellinfo.SellList, sell)
	}

	return resprofilesellinfo, nil

}

func (m *mariadbHandler) ReadProfileBuyInfo(userid int64) (*ResProfileBuyInfo, error) {
	var resprofilebuyinfo = &ResProfileBuyInfo{}
	var buy Buy

	stmt, err := m.db.Prepare(`SELECT 
						p.id, 
						p.title,
						bh.created_at,
						bh.seller_id,
						bh.seller_nickname,
						bh.price
						FROM user AS u 
						INNER JOIN buy_history as bh ON bh.user_id = u.id
						INNER JOIN project AS p ON p.id = bh.project_id
						WHERE u.id = ?`)
	if err != nil {
		log.Println("[ERR] prepare stmt err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(userid)
	if err != nil {
		log.Println("[ERR] stmt query err : ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&buy.ID, &buy.Title, &buy.Date, &buy.SellerID, &buy.SellerNickName, &buy.Price)
		if err != nil {
			log.Println("[ERR] stmt query err : ", err)
			return nil, err
		}
		resprofilebuyinfo.BuyList = append(resprofilebuyinfo.BuyList, buy)
	}

	return resprofilebuyinfo, nil

}

func (m *mariadbHandler) ReadProfileWithdrawInfo(userid int64) (*ResProfileWithdrawInfo, error) {
	var resprofilewithdrawinfo = &ResProfileWithdrawInfo{}
	var withdraw Withdraw

	stmt, err := m.db.Prepare(`SELECT 
				  u.id, 
				  wh.requested_at,
				  wh.completed_at,
				  wh.amount
				  FROM user AS u  
				  INNER JOIN withdraw_history as wh ON wh.user_id = u.id
				  WHERE u.id= ?`)
	if err != nil {
		log.Println("[ERR] prepare stmt err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(userid)
	if err != nil {
		log.Println("[ERR] stmt query err : ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&withdraw.ID, &withdraw.RequestedDate, &withdraw.CompleteDate, &withdraw.Amount)
		if err != nil {
			log.Println("[ERR] stmt query err : ", err)
			return nil, err
		}
		resprofilewithdrawinfo.WithdrawList = append(resprofilewithdrawinfo.WithdrawList, withdraw)
	}

	return resprofilewithdrawinfo, nil

}

func (m *mariadbHandler) UpdateModificationUserInfo(reqmodinfo *ReqModificationUserInfo) error {
	tx, err := m.db.Begin()
	if err != nil {
		log.Println("[ERR] begin err : ", err)
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE user SET name=?,nickname=?,email=?,agree_email_marketing=?, introduction = ?, image_link=? WHERE user.id =?`)
	if err != nil {
		log.Println("[ERR] Prepare statement err : ", err)
		return err
	}

	defer stmt.Close()

	stmt.Exec()

	result, err := stmt.Exec(reqmodinfo.Name, reqmodinfo.Nickname, reqmodinfo.Email, reqmodinfo.AgreeEmailMarketing, reqmodinfo.Introduction, reqmodinfo.ImageLink, reqmodinfo.UserID)
	if err != nil {
		log.Println("[ERR] Exec err : ", err)
		return err
	}

	rowcnt, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERR] rows affected err : ", err)
		return err
	}

	log.Println("affected rwos count : ", rowcnt)

	err = tx.Commit()
	if err != nil {
		log.Println("[ERR] commit err : ", err)
		return err
	}

	return nil
}

func (m *mariadbHandler) ReadModificationUserInfo(userid int64) (*ResModificationUserInfo, error) {
	var resmodificationuserinfo = &ResModificationUserInfo{}

	stmt, err := m.db.Prepare(`SELECT 
	name, 
	nickname,
	email,
	agree_email_marketing,
	introduction,
	image_link
	FROM user   
	WHERE user.id= ?`)

	if err != nil {
		log.Println("[ERR] prepare statement err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(userid)
	if err != nil {
		log.Println("[ERR] stmt query err : ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&resmodificationuserinfo.Name, &resmodificationuserinfo.Nickname, &resmodificationuserinfo.Email, &resmodificationuserinfo.AgreeEmailMarketing, &resmodificationuserinfo.Introduction, &resmodificationuserinfo.ImageLink)
		if err != nil {
			log.Println("[ERR] rows scan err : ", err)
			return nil, err
		}
	}

	return resmodificationuserinfo, err
}

func (m *mariadbHandler) ReadProfileArtistInfo(artistid int64) (*ResArtistProfileInfo, error) {
	var resartistinfo = &ResArtistProfileInfo{}
	var project Project

	stmt, err := m.db.Prepare(`SELECT 
	u.id,
	u.nickname,
	u.introduction,
	(SELECT COUNT(id) FROM project WHERE user_id = ?),
	(SELECT COUNT(id) FROM sell_history WHERE user_id = ?),
	u.image_link,
	p.id,
	p.title,
	p.description,
	i.link,
	p.created_at,
	p.sell_count,
	u.id, 
	p.comment_count,
	p.total_upvote_count,
	p.price,
	p.beta
	FROM user AS u 
	INNER JOIN project AS p ON p.user_id = u.id
	INNER JOIN image AS i ON i.project_id = p.id
	INNER JOIN 
	(SELECT project_id , MIN(created_at) created_at 
	FROM image GROUP BY project_id) AS ii ON ii.project_id = i.project_id AND ii.created_at = i.created_at
	WHERE u.id = ? GROUP BY p.id
`)

	if err != nil {
		log.Println("[ERR] prepare statement err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(artistid, artistid, artistid)
	if err != nil {
		log.Println("[ERR] query err : ", err)
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&resartistinfo.ArtistInfo.ArtistID, &resartistinfo.ArtistInfo.ArtistNickName, &resartistinfo.ArtistInfo.Introduction, &resartistinfo.ArtistInfo.ProjectCount, &resartistinfo.ArtistInfo.SellCount, &resartistinfo.ArtistInfo.ImageLink, &project.ProjectID, &project.Title, &project.Desc, &project.ImageLink, &project.CreatedAt, &project.SellCount, &project.UserNickName, &project.CommontCount, &project.UpvoteCount, &project.Price, &project.Beta)
		if err != nil {
			log.Println("[ERR] scan err : ", err)
			return nil, err
		}

		resartistinfo.ProjectList = append(resartistinfo.ProjectList, project)

	}

	return resartistinfo, nil
}

func (m *mariadbHandler) ReadPersonalInformation(userid int64) (*ResPersonalInformation, error) {
	var respersonalinformation = &ResPersonalInformation{}

	stmt, err := m.db.Prepare(`SELECT 
	bank, 
	account,
	agree_policy
	FROM user  
	WHERE user.id= ?`)

	if err != nil {
		log.Println("[ERR] prepare statement err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(userid)
	if err != nil {
		log.Println("[ERR] stmt query err : ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&respersonalinformation.Account.Bank, &respersonalinformation.Account.Account, &respersonalinformation.Account.AgreePolicy)
		if err != nil {
			log.Println("[ERR] rows scan err : ", err)
			return nil, err
		}
	}

	return respersonalinformation, err

}

func (m *mariadbHandler) UpdatePersonalInformation(reqpersonalinformation *ReqPersonalInformation) error {

	tx, err := m.db.Begin()
	if err != nil {
		log.Println("[ERR] begin err : ", err)
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE user 
				SET 
				bank = ?, 
				account = ?,
				agree_policy = ? ,
				updated_at = NOW()
				WHERE user.id = ?`)
	if err != nil {
		log.Println("[ERR] prepare statement err : ", err)
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(reqpersonalinformation.Account.Bank, reqpersonalinformation.Account.Account, reqpersonalinformation.Account.AgreePolicy, reqpersonalinformation.Account.UserID)
	if err != nil {
		log.Println("[ERR] statment execution err : ", err)
		return err
	}

	rowcnt, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERR] rows affected err : ", err)
		return err
	}

	log.Println("affected rwos count : ", rowcnt)

	err = tx.Commit()
	if err != nil {
		log.Println("[ERR] commit err : ", err)
		return err
	}

	return nil
}
