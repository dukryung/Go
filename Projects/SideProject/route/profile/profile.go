package profile

import (
	"database/sql"
	"log"
	"net/http"
	"sideproject/route/user"

	"github.com/gin-gonic/gin"
)

type Profile struct {
	DB *sql.DB
}

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

type Account struct {
	UserID      int64  `json:"user_id"`
	Bank        int    `json:"bank"`
	Account     string `json:"account"`
	AgreePolicy bool   `json:"agree_policy"`
}

func (pf *Profile) Routes(route *gin.RouterGroup) {

	gru := route.Group("/user")
	{
		gru.GET("/index", pf.getProfile)
		gru.GET("/frame", pf.getProfileFrameInfo)
		gru.GET("/project", pf.getProfileProject)
		gru.GET("/sell", pf.getProfileSell)
		gru.GET("/buy", pf.GetProfileBuy)
		gru.GET("/withdraw", pf.getProfileWithdraw)
		gru.GET("/modification", pf.getModificationUserInfo)
		gru.PUT("/modification", pf.putModificationUserInfo)
	}

	grpl := route.Group("/personal")
	{
		grpl.GET("/index", pf.getPersonalIndex)
		grpl.GET("/information", pf.getPersonalInformation)
		grpl.PUT("/information", pf.putPersonalInformation)
	}
	grart := route.Group("/artist")
	{

		grart.GET("/index", pf.getProfileArtist)
		//------ {TODO : project and frame split
		grart.GET("/information", pf.getProfileArtistInfo)
		//------ TODO : project and frame split }
	}

}

func (pf *Profile) getProfile(c *gin.Context) {
	c.HTML(http.StatusOK, "profileuser.html", gin.H{
		"title": "Profile User Page",
	})
}
func (pf *Profile) getProfileFrameInfo(c *gin.Context) {
	userinfo := &user.JSUser{}

	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user idd err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	resprofileframeinfo, err := pf.ReadProfileFrameInfo(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read profile frame information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofileframeinfo)
}
func (pf *Profile) getProfileProject(c *gin.Context) {
	userinfo := &user.JSUser{}

	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user idd err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	resprofileprojectinfo, err := pf.ReadProfileProjectInfo(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read profile project information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofileprojectinfo)
}

func (pf *Profile) getProfileSell(c *gin.Context) {
	userinfo := &user.JSUser{}

	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	resprofilesellinfo, err := pf.ReadProfileSellInfo(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read profile sell history information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofilesellinfo)
}

func (pf *Profile) GetProfileBuy(c *gin.Context) {
	userinfo := &user.JSUser{}

	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	resprofilebuyinfo, err := pf.ReadProfileBuyInfo(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read profile buy history information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofilebuyinfo)
}

func (pf *Profile) getProfileWithdraw(c *gin.Context) {
	userinfo := &user.JSUser{}

	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	resprofilebuyinfo, err := pf.ReadProfileWithdrawInfo(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read profile withdraw history information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofilebuyinfo)
}

//GetModificationUserInfoHandler is function to get user profile infomation.
func (pf *Profile) getModificationUserInfo(c *gin.Context) {
	userinfo := &user.JSUser{}

	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	resmodificationuserinfo, err := pf.ReadModificationUserInfo(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read modification user info err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, resmodificationuserinfo)
}

//PutModificationUserInfoHandler is fuction to edit user profile information.
func (pf *Profile) putModificationUserInfo(c *gin.Context) {
	var reqmodinfo = &ReqModificationUserInfo{}
	err := c.ShouldBindJSON(reqmodinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	err = pf.UpdateModificationUserInfo(reqmodinfo)
	if err != nil {
		log.Println("[ERR] update user info ")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, nil)

}

func (pf *Profile) getPersonalIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "personalinformation.html", gin.H{
		"title": "Personal Information Page",
	})
}
func (pf *Profile) getPersonalInformation(c *gin.Context) {
	var userinfo = &user.JSUser{}
	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	respersonalinfomation, err := pf.ReadPersonalInformation(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read personal information err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, respersonalinfomation)
}

func (pf *Profile) putPersonalInformation(c *gin.Context) {
	var reqpersonalinformation = &ReqPersonalInformation{}
	err := c.ShouldBindJSON(reqpersonalinformation)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	err = pf.UpdatePersonalInformation(reqpersonalinformation)
	if err != nil {
		log.Println("[ERR] failed to update personal information err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (pf *Profile) getProfileArtist(c *gin.Context) {
	c.HTML(http.StatusOK, "profileartist.html", gin.H{
		"title": "profile artist  Page",
	})
}

func (pf *Profile) ReadProfileFrameInfo(userid int64) (*ResProfileFrameInfo, error) {

	var resprofileframeinfo = &ResProfileFrameInfo{}

	stmt, err := pf.DB.Prepare(`SELECT 
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

func (pf *Profile) getProfileArtistInfo(c *gin.Context) {
	var artistprofileinfo = &ArtistProfile{}
	err := c.ShouldBindJSON(artistprofileinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	resartistinfo, err := pf.ReadProfileArtistInfo(artistprofileinfo.ArtistID)
	if err != nil {
		log.Println("[ERR] failed to read artist information err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resartistinfo)
}

func (pf *Profile) ReadProfileProjectInfo(userid int64) (*ResProfileProjectInfo, error) {
	var resprofileprojectinfo = &ResProfileProjectInfo{}
	var project Project

	stmt, err := pf.DB.Prepare(`SELECT 
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

func (pf *Profile) ReadProfileSellInfo(userid int64) (*ResProfileSellInfo, error) {
	var resprofilesellinfo = &ResProfileSellInfo{}
	var sell Sell

	stmt, err := pf.DB.Prepare(`  SELECT 
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

func (pf *Profile) ReadProfileBuyInfo(userid int64) (*ResProfileBuyInfo, error) {
	var resprofilebuyinfo = &ResProfileBuyInfo{}
	var buy Buy

	stmt, err := pf.DB.Prepare(`SELECT 
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

func (pf *Profile) ReadProfileWithdrawInfo(userid int64) (*ResProfileWithdrawInfo, error) {
	var resprofilewithdrawinfo = &ResProfileWithdrawInfo{}
	var withdraw Withdraw

	stmt, err := pf.DB.Prepare(`SELECT 
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

func (pf *Profile) UpdateModificationUserInfo(reqmodinfo *ReqModificationUserInfo) error {
	tx, err := pf.DB.Begin()
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

func (pf *Profile) ReadModificationUserInfo(userid int64) (*ResModificationUserInfo, error) {
	var resmodificationuserinfo = &ResModificationUserInfo{}

	stmt, err := pf.DB.Prepare(`SELECT 
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

func (pf *Profile) ReadProfileArtistInfo(artistid int64) (*ResArtistProfileInfo, error) {
	var resartistinfo = &ResArtistProfileInfo{}
	var project Project

	stmt, err := pf.DB.Prepare(`SELECT 
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

func (pf *Profile) ReadPersonalInformation(userid int64) (*ResPersonalInformation, error) {
	var respersonalinformation = &ResPersonalInformation{}

	stmt, err := pf.DB.Prepare(`SELECT 
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

func (pf *Profile) UpdatePersonalInformation(reqpersonalinformation *ReqPersonalInformation) error {

	tx, err := pf.DB.Begin()
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
