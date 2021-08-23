package iamport

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// Client need for networking with iamport's server.
type Client struct {
	APIKey      string
	APISecret   string
	AccessToken accessToken
	HTTP        *http.Client
}

type accessToken struct {
	Token   string
	Expired time.Time
}

type Iamport struct {
	DB *sql.DB
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
	Iid string `json:"imp_uid"`
	Mid string `json:"merchant_uid"`
}

func (iam *Iamport) Routes(route *gin.RouterGroup) {
	route.GET("/payment/information", iam.getPAYMT)
	route.POST("/payment/complete", iam.checkPAYMT)
}

//PAYT is abbreviation of payment.
//checkPAYMT checks the payment information at iamport.
func (iam *Iamport) checkPAYMT(c *gin.Context) {
	reqrsc := &REQRSC{}
	err := c.ShouldBindJSON(reqrsc)
	if err != nil {
		log.Println("[ERR] failed to extract request resourece err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}
}

//IMPT is abbreviation of iamport.
//GetIMPTToken get authentication token from iamport.
func (cli *Client) GetIMPTToken() error {
	if cli.APIKey == "" {
		return errors.New("iamport: APIKey is missing")
	}

	if cli.APISecret == "" {
		return errors.New("iamport: APISecret is missing")
	}

	form := url.Values{}
	form.Set("imp_key", cli.APIKey)
	form.Set("imp_secret", cli.APISecret)

	res, err := cli.HTTP.PostForm("https://api.iamport.kr/users/getToken", form)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return errors.New("iamport: unauthorized")
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("iamport: unknown error")
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
		return fmt.Errorf("iamport: %s", data.Message)
	}

	cli.AccessToken.Token = data.Response.AccessToken
	cli.AccessToken.Expired = time.Unix(data.Response.ExpiredAt, 0)

	return nil
}

//PAYT is abbreviation of payment.
//getPAYT reads information required for payment request.
func (iam *Iamport) getPAYMT(c *gin.Context) {
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
