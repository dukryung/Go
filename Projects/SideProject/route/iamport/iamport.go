package iamport

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

var APIKey = "5159389833177379"
var APISecret = "GkmriPyRX9Qf0O8m7i4TGOThaYgDBPpi9MKsCikzsaPtZXQ7gU1XJUeTTvwukWMuIwWkPgUirJLrGByX"

type accessToken struct {
	Token   string
	Expired time.Time
}

type Iamport struct {
	DB          *sql.DB
	APIKey      string
	APISecret   string
	AccessToken accessToken
	HTTP        *http.Client
}

// Payment information
type PAYMT struct {
	ImpUID        string `json:"imp_uid"`
	MerchantUID   string `json:"merchant_uid"`
	PayMethod     string `json:"pay_method"`
	PGProvider    string `json:"pg_provider"`
	PGTID         string `json:"pg_tid"`
	ApplyNum      string `json:"apply_num"`
	CardName      string `json:"card_name"`
	CardQuota     int    `json:"card_quota"`
	VBankName     string `json:"vbank_name"`
	VBankNum      string `json:"vbank_num"`
	VBankHolder   string `json:"vbank_holder"`
	Name          string `json:"name"`
	Amount        int64  `json:"amount"`
	CancelAmount  int64  `json:"cancel_amount"`
	BuyerName     string `json:"buyer_name"`
	BuyerEmail    string `json:"buyer_email"`
	BuyerTel      string `json:"buyer_tel"`
	BuyerAddr     string `json:"buyer_addr"`
	BuyerPostCode string `json:"buyer_postcode"`
	CustomData    string `json:"custom_data"`
	UserAgent     string `json:"user_agent"`
	Status        string `json:"status"`
	PaidAt        int64  `json:"paid_at"`
	FailedAt      int64  `json:"failed_at"`
	CanceledAt    int64  `json:"canceled_at"`
	FailReason    string `json:"fail_reason"`
	CancelReason  string `json:"cancel_reason"`
	ReceiptURL    string `json:"receipt_url"`
}

//REQRSC is abbreviation of request  resource.
type REQRSC struct {
	Uid int64  `json:"user_id"`
	Pid int64  `json:"project_id"`
	Aid int64  `json:"artist_id"`
	Iid string `json:"imp_uid"`
	Mid string `json:"merchant_uid"`
}

func (iam *Iamport) Routes(route *gin.RouterGroup) {
	route.GET("/payment/information", iam.GetPAYMT)
	route.POST("/payment/complete", iam.CheckPAYMT)
	//---{TODO : cancel logic
	route.POST("/payment/cancel")
	//---}
}

func (iam *Iamport) CancelPAYMT(c *gin.Context) {

}

//PAYT is abbreviation of payment.
//getPAYT reads information required for payment request.
func (iam *Iamport) GetPAYMT(c *gin.Context) {
	reqrsc := &REQRSC{}
	err := c.ShouldBindJSON(reqrsc)
	if err != nil {
		log.Println("[ERR] failed to extract request resourece err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	var p = &PAYMT{}
	err = iam.getPAYMTRSCTODB(p)
	if err != nil {
		log.Println("[ERR] failed to extract payment resource err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	err = iam.getBUYRSCTODB(p, reqrsc.Uid)
	if err != nil {
		log.Println("[ERR] failed to extract buyer resource err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	err = iam.getPJRSCTODB(p, reqrsc.Pid)
	if err != nil {
		log.Println("[ERR] failed to extract merchant resource err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	err = getMerchantid(p)
	if err != nil {
		log.Println("[ERR] failed to extract merchant id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, p)
}

//PAYT is abbreviation of payment.
//checkPAYMT checks the payment information at iamport.
func (iam *Iamport) CheckPAYMT(c *gin.Context) {
	reqrsc := &REQRSC{}
	err := c.ShouldBindJSON(reqrsc)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	fmt.Println("request resource : ", reqrsc)

	err = iam.GetIMPTToken()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	pay, err := iam.GetPAYMTImpUID(reqrsc.Iid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	fmt.Println("pay:", pay)

	p, err := iam.GetPrice(reqrsc.Pid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	err = iam.ComparePrice(pay.Amount, p)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	err = iam.SavePAYMTTODB(pay, reqrsc)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, nil)

}
func (iam *Iamport) SavePAYMTTODB(pay PAYMT, reqsrc *REQRSC) error {

	tx, err := iam.DB.Begin()
	if err != nil {
		return errors.New("[ERR] fail transaction begin")
	}
	defer tx.Rollback()

	err = insertBuyHistory(tx, pay, reqsrc)
	if err != nil {
		return err
	}

	err = insertSellHistory(tx, pay, reqsrc)
	if err != nil {
		return err
	}

	err = insertSaleList(tx, reqsrc)
	if err != nil {
		return err
	}

	err = updatePAYMTConcerned(tx, pay, reqsrc)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

//updatePAYMTConcerned updates project sell_count and cash
func updatePAYMTConcerned(tx *sql.Tx, pay PAYMT, reqsrc *REQRSC) error {

	stmt, err := tx.Prepare(`UPDATE user AS u JOIN project AS p ON u.id = p.user_id SET cash= cash + ?, sell_count = sell_count + 1 WHERE p.id = ?`)
	if err != nil {
		return errors.New("[ERR] failed to prepare statement in updatePAYMTConcern")
	}
	defer stmt.Close()

	_, err = stmt.Exec(pay.Amount, reqsrc.Pid)
	if err != nil {
		return errors.New("[ERR] failed to execute insert sale information in InsertSaleList")
	}

	return nil
}

//insertSaleList inserts projects sales.
func insertSaleList(tx *sql.Tx, reqsrc *REQRSC) error {
	stmt, err := tx.Prepare(`INSERT INTO salelist (user_id, project_id) 
	VALUES(?,?)`)
	if err != nil {
		return errors.New("[ERR] failed  to prepare statement in InsertSaleList")
	}
	defer stmt.Close()

	_, err = stmt.Exec(reqsrc.Uid, reqsrc.Pid)
	if err != nil {
		return errors.New("[ERR] failed to execute insert sale information in InsertSaleList")
	}

	return nil
}

//SBRSC is abbreviation of Seller and buyer resource
//SelectSBRSC reads seller and buyer information
func selectSBRSC(tx *sql.Tx, userid int64) (string, error) {
	stmt, err := tx.Prepare(`SELECT nickname FROM user WHERE id = ? `)
	if err != nil {
		return "", errors.New("[ERR] failed to prepare statement in SelectBuyerRSC")
	}

	defer stmt.Close()

	rows, err := stmt.Query(userid)
	if err != nil {
		return "", errors.New("[ERR] failed to query in SelectBuyerRSC")
	}
	defer rows.Close()

	var nickname string
	for rows.Next() {
		err = rows.Scan(&nickname)
		if err != nil {
			return "", errors.New("[ERR] failed to scan in SelectBuyerRSC")
		}
	}

	return nickname, nil
}

//InsertBuyHistory saves buy history to DB.
func insertBuyHistory(tx *sql.Tx, pay PAYMT, reqsrc *REQRSC) error {

	snickname, err := selectSBRSC(tx, reqsrc.Aid)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO buy_history 
	(user_id, project_id, seller_id, seller_nickname, price,impuid, merchant_id)
	VALUES(?,?,?,?,?,?,?)`)
	if err != nil {
		return errors.New("[ERR] failed to prepare statement in InsertBuyHistory")
	}
	defer stmt.Close()

	_, err = stmt.Exec(reqsrc.Uid, reqsrc.Pid, reqsrc.Aid, snickname, pay.Amount, pay.ImpUID, pay.MerchantUID)
	if err != nil {
		return errors.New("[ERR] failed to execute statement")
	}

	return nil
}

//InsertSellHistory saves buy history to DB.
func insertSellHistory(tx *sql.Tx, pay PAYMT, reqsrc *REQRSC) error {

	bnickname, err := selectSBRSC(tx, reqsrc.Uid)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO buy_history 
	(user_id, project_id, seller_id, seller_nickname, price,impuid, merchant_id)
	VALUES(?,?,?,?,?,?,?)`)
	if err != nil {
		return errors.New("[ERR] failed to prepare statement in InsertBuyHistory")
	}
	defer stmt.Close()

	_, err = stmt.Exec(reqsrc.Aid, reqsrc.Pid, reqsrc.Uid, bnickname, pay.Amount, pay.ImpUID, pay.MerchantUID)
	if err != nil {
		return errors.New("[ERR] failed to execute statement")
	}

	return nil
}

func (iam *Iamport) ComparePrice(amount int64, price int64) error {

	if price != amount {
		return errors.New("[ERR] not same price ")
	}

	return nil
}
func (iam *Iamport) GetPrice(projectid int64) (int64, error) {
	stmt, err := iam.DB.Prepare(`SELECT price FROM project WHERE id = ?`)
	if err != nil {

		return 0, errors.New("[ERR] failed to prepare statement project price query")
	}

	defer stmt.Close()

	rows, err := stmt.Query(projectid)
	if err != nil {
		return 0, errors.New("[ERR] failed to query err")
	}

	var price int64
	for rows.Next() {
		err = rows.Scan(price)
		if err != nil {
			return 0, errors.New("[ERR] failed to query err")
		}
	}

	return price, nil
}

//IMPT is abbreviation of iamport.
//GetIMPTToken get authentication token from iamport.
func (iam *Iamport) GetIMPTToken() error {
	reqsrc := struct {
		ImpKey    string `json:"imp_key"`
		ImpSecret string `json:"imp_secret"`
	}{}

	reqsrc.ImpKey = iam.APIKey
	reqsrc.ImpSecret = iam.APISecret

	if iam.APIKey == "" {
		return errors.New("[ERR] APIKey is missing")
	}

	if iam.APISecret == "" {
		return errors.New("[ERR] APISecret is missing")
	}

	imdata, err := json.Marshal(reqsrc)
	if err != nil {
		return errors.New("[ERR] failed to json marshal")
	}

	buf := bytes.NewBuffer(imdata)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.iamport.kr/users/getToken", buf)
	if err != nil {
		return errors.New("[ERR] failed to json marshal")
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return errors.New("[ERR] client do err")
	}

	if res.StatusCode == http.StatusUnauthorized {
		return errors.New("[ERR] unauthorized")
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("[ERR] unknown error")
	}

	data := struct {
		Code     int    `json:"code"`
		Message  string `json:"message"`
		Response struct {
			AccessToken string `json:"access_token"`
			ExpiredAt   int64  `json:"expired_at"`
			Now         int64  `json:"now"`
		} `json:"response"`
	}{}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return err
	}

	if data.Code != 0 {
		return fmt.Errorf("[ERR] %s", data.Message)
	}

	iam.AccessToken.Token = data.Response.AccessToken
	iam.AccessToken.Expired = time.Unix(data.Response.ExpiredAt, 0)

	return nil
}

func (iam *Iamport) authorization() (string, error) {
	now := time.Now()

	switch {
	case iam.AccessToken.Token == "",
		iam.AccessToken.Expired.IsZero(),
		!iam.AccessToken.Expired.Before(now):

		err := iam.GetIMPTToken()
		if err != nil {
			return "", err
		}
	}

	return iam.AccessToken.Token, nil
}

// GetPAYMTImpUID reads  payment information using iamport number.
// GET /payments/{imp_uid}
func (iam *Iamport) GetPAYMTImpUID(impuid string) (PAYMT, error) {
	data := struct {
		Code     int    `json:"code"`
		Message  string `json:"message"`
		Response PAYMT  `json:"response"`
	}{}

	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.iamport.kr/payments/%s", impuid), nil)
	if err != nil {
		return data.Response, err
	}

	auth, err := iam.authorization()
	if err != nil {
		return data.Response, err
	}

	req.Header.Set("Authorization", auth)

	res, err := client.Do(req)
	if err != nil {
		return data.Response, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return data.Response, errors.New("[ERR] unauthorized")
	}

	if res.StatusCode == http.StatusNotFound {
		return data.Response, errors.New("[ERR] invalid imp_uid")
	}

	if res.StatusCode != http.StatusOK {
		return data.Response, errors.New("[ERR] unknown error")
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return data.Response, err
	}

	if data.Code != 0 {
		return data.Response, fmt.Errorf("[ERR] %s", data.Message)
	}

	return data.Response, nil
}

//PAYMTRSC is abbreviation of payment resource.
//getPAYMTRSC reads pg and payment method and iamport uid in DB.
func (iam *Iamport) getPAYMTRSCTODB(p *PAYMT) error {
	stmt, err := iam.DB.Prepare(`SELECT impuid, pg, pay_method FROM payment`)
	if err != nil {
		log.Println("[ERR] prepare stmt err : ", err)
		return nil
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Println("[ERR] stmt query err : ", err)
		return nil
	}

	for rows.Next() {
		err = rows.Scan(p.ImpUID, p.PGProvider, p.PayMethod)
		if err != nil {
			return nil
		}
	}

	return nil

}

//BUYRSC is abbreviation of buyer resource.
//getBUYRSCTODB reads buyer information in DB.
func (iam *Iamport) getBUYRSCTODB(p *PAYMT, userid int64) error {
	stmt, err := iam.DB.Prepare(`SELECT name, email, phonenum FROM user WHERE id = ?`)
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

	for rows.Next() {
		err = rows.Scan(p.BuyerName, p.BuyerEmail, p.BuyerTel)
		if err != nil {
			return nil
		}
	}

	return nil
}

//PJRSC is abbreviation of project resource.
//getPJRSCTODB reads project information in DB.
func (iam *Iamport) getPJRSCTODB(p *PAYMT, projectid int64) error {
	stmt, err := iam.DB.Prepare(`SELECT title, price FROM project WHERE id = ?`)
	if err != nil {
		log.Println("[ERR] prepare stmt err : ", err)
		return nil
	}

	defer stmt.Close()

	rows, err := stmt.Query(projectid)
	if err != nil {
		log.Println("[ERR] stmt query err : ", err)
		return nil
	}

	for rows.Next() {
		err = rows.Scan(p.Name, p.Amount)
		if err != nil {
			return nil
		}
	}

	return nil
}

func getMerchantid(p *PAYMT) error {
	n := getnowtime()
	u, err := getuuid()
	if err != nil {
		fmt.Println("[ERR] failed getuuid")
		return err
	}

	p.MerchantUID = makeMerchantid(n, u)
	return nil
}

func getnowtime() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}

func getuuid() (string, error) {
	gen := uuid.NewGen()
	u, err := gen.NewV4()
	if err != nil {
		fmt.Println("[ERR] failed generate new v4 uuid")
		return "", err
	}
	return u.String(), nil
}

func makeMerchantid(now string, uuid string) string {
	return now + "_" + uuid
}
