package common

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2/kakao"

	"golang.org/x/oauth2/facebook"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Social        string `json:"social"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

type AuthKakaoUserInfo struct {
	ID           int              `json:"id"`
	KakaoAccount KakaoAccountInfo `json:"kakao_account"`
}
type KakaoAccountInfo struct {
	Email string `json:"email"`
}
type AuthNaverUserInfo struct {
	NaverAccount NaverAccountInfo `json:"response"`
}
type NaverAccountInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

var googleOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     "10397292751-jqmb2g27a9pks1qba5m4h576qacss16b.apps.googleusercontent.com",
	ClientSecret: "m65Xt4QPuloEz3bgIBfF2GbG",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func (p *project) googleLoginHandler(c *gin.Context) {
	state := generateStateOauthCookie(c)
	url := googleOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func generateStateOauthCookie(c *gin.Context) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(c.Writer, cookie)

	return state
}

var store = sessions.NewCookieStore([]byte("a12574e0-a68b-4e99-82ba-2f1d950b8126"))

func (p *project) googleAuthCallback(c *gin.Context) {

	oauthstate, err := c.Request.Cookie("oauthstate")
	if err != nil {
		log.Println("[ERR] read Cookie err :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	if c.Request.FormValue("state") != oauthstate.Value {
		log.Printf("[ERR] do not match cookie : %s , state : %s  ", oauthstate.Value, c.Request.FormValue("state"))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	data, err := getGoogleUserInfo(c.Request.FormValue("code"))
	if err != nil {
		log.Println("[ERR] get user info err : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	var userinfo AuthUserInfo
	err = json.Unmarshal(data, &userinfo)
	if err != nil {
		log.Println("[ERR] unmarshal err : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	session, err := store.Get(c.Request, "session")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	session.Values["id"] = userinfo.ID
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	userinfo.Social = "google"
	err = p.db.CreateUserInfo(userinfo)
	if err != nil {
		log.Println("[ERR] session save err :", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	var resemailinfo ResEmailInfo
	resemailinfo.Email = userinfo.Email

	c.JSON(http.StatusOK, resemailinfo)
}

//CheckSessionValidity is function to check whether the ssesion is valid or not.
func CheckSessionValidity(c *gin.Context) {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		log.Println("[ERR] get sesion err : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		c.Abort()
		return
	}

	val := session.Values["id"]
	//test -----------------
	val = 1
	//-------------------test
	if val == nil {
		log.Println("[ERR] session value err : ", val)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		c.Abort()
		return
	}

	c.Next()
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	return ioutil.ReadAll(resp.Body)

}

var facebookOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/facebook/callback",
	ClientID:     "983065818893651",
	ClientSecret: "025390e6599ac4c9afabf0b7a25040f6",
	Scopes:       []string{"public_profile", "email"},
	Endpoint:     facebook.Endpoint,
}

func (p *project) facebookLoginHandler(c *gin.Context) {
	state := generateStateOauthCookie(c)
	url := facebookOauthConfig.AuthCodeURL(state)
	fmt.Println("url:", url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (p *project) facebookAuthCallback(c *gin.Context) {

	oauthstate, err := c.Request.Cookie("oauthstate")
	if err != nil {
		log.Println("r Cookie err :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	if c.Request.FormValue("state") != oauthstate.Value {
		log.Printf("[ERR] do not match cookie : %s , state : %s  ", oauthstate.Value, c.Request.FormValue("state"))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	data, err := getfacebookUserInfo(c.Request.FormValue("code"))
	if err != nil {
		log.Println("[ERR] get facebook user info err :", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	log.Println("data ", string(data))

	var userinfo AuthUserInfo
	err = json.Unmarshal(data, &userinfo)
	if err != nil {
		log.Println("[ERR] unmarshal err : ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	session, err := store.Get(c.Request, "session")
	if err != nil {
		log.Println("[ERR] store get err :", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	session.Values["id"] = userinfo.ID
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Println("[ERR] session save err :", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	userinfo.Social = "facebook"
	err = p.db.CreateUserInfo(userinfo)
	if err != nil {
		log.Println("[ERR] session save err :", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	var resemailinfo ResEmailInfo
	resemailinfo.Email = userinfo.Email

	c.JSON(http.StatusOK, resemailinfo)

}

const oauthFacebookUrlAPI = "https://graph.facebook.com/me?fields=id,name,email&access_token="

func getfacebookUserInfo(code string) ([]byte, error) {
	token, err := facebookOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	resp, err := http.Get(oauthFacebookUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	return ioutil.ReadAll(resp.Body)

}

var kakaoOauthConfig = oauth2.Config{
	ClientID:     "7b23cc14ba23a32a2b9cb16fde1b92b1",
	ClientSecret: "G92PaUhaNWEeMLTBsIvfHU6vmk6PiUzy",
	RedirectURL:  "http://localhost:8080/auth/kakao/callback",
	//Scopes:       []string{"profile"},
	Endpoint: kakao.Endpoint,
}

func (p *project) kakaoLoginHandler(c *gin.Context) {
	state := generateStateOauthCookie(c)
	url := kakaoOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (p *project) kakaoAuthCallback(c *gin.Context) {

	oauthstate, err := c.Request.Cookie("oauthstate")
	if err != nil {
		log.Println("[ERR] r.Cookie err : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	if c.Request.FormValue("state") != oauthstate.Value {
		log.Printf("[ERR] do not match cookie : %s , state : %s  ", oauthstate.Value, c.Request.FormValue("state"))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	data, err := getKakaoUserInfo(c.Request.FormValue("code"))
	if err != nil {
		log.Println("[ERR] err  in getKakaoUserInfo function : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})

	}

	var kakaouserinfo AuthKakaoUserInfo
	err = json.Unmarshal(data, &kakaouserinfo)
	if err != nil {
		log.Println("[ERR] unmarshal err : ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	var userinfo AuthUserInfo

	userinfo.ID = strconv.Itoa(kakaouserinfo.ID)
	userinfo.Email = kakaouserinfo.KakaoAccount.Email

	session, err := store.Get(c.Request, "session")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	session.Values["id"] = userinfo.ID

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	userinfo.Social = "kakao"
	err = p.db.CreateUserInfo(userinfo)
	if err != nil {
		log.Println("[ERR] session save err :", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	var resemailinfo ResEmailInfo
	resemailinfo.Email = userinfo.Email

	c.JSON(http.StatusOK, resemailinfo)

}

const oauthKakaoUrlAPI = "https://kapi.kakao.com/v2/user/me?access_token="

func getKakaoUserInfo(code string) ([]byte, error) {
	token, err := kakaoOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	res, err := http.Get(oauthKakaoUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll error  %s\n: ", err.Error())
	}

	return data, err
}

var naverOauthConfig = oauth2.Config{
	ClientID:     "ILDfZE0DKDPj1rBk9sfJ",
	ClientSecret: "zeBbkQLBp5",
	RedirectURL:  "http://localhost:8080/auth/naver/callback",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://nid.naver.com/oauth2.0/authorize",
		TokenURL: "https://nid.naver.com/oauth2.0/token",
	},
}

func (p *project) naverLoginHndler(c *gin.Context) {
	state := generateStateOauthCookie(c)
	url := naverOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusOK, url)
}

func (p *project) naverAuthCallback(c *gin.Context) {

	oauthstate, err := c.Request.Cookie("oauthstate")
	if err != nil {
		log.Println("[ERR] r.Cookie err : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	if c.Request.FormValue("state") != oauthstate.Value {
		log.Printf("[ERR] do not match cookie : %s , state : %s  ", oauthstate.Value, c.Request.FormValue("state"))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	data, err := getNaverUserInfo(c.Request.FormValue("code"))
	if err != nil {
		log.Println("[ERR] err  in getKakaoUserInfo function : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	var naveruserinfo AuthNaverUserInfo

	err = json.Unmarshal(data, &naveruserinfo)
	if err != nil {
		log.Println(" [ERR] Unmarshal err : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	var userinfo AuthUserInfo

	userinfo.ID = naveruserinfo.NaverAccount.ID
	userinfo.Email = naveruserinfo.NaverAccount.Email

	session, err := store.Get(c.Request, "session")
	if err != nil {
		log.Println("[ERR] store Get err : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	session.Values["id"] = userinfo.ID

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Println("[ERR] session save err : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	userinfo.Social = "naver"
	err = p.db.CreateUserInfo(userinfo)
	if err != nil {
		log.Println("[ERR] create  user info err : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	var resemailinfo ResEmailInfo
	resemailinfo.Email = userinfo.Email

	c.JSON(http.StatusOK, resemailinfo)

}

const oauthNaverUrlAPI = "https://openapi.naver.com/v1/nid/me?access_token="

func getNaverUserInfo(code string) ([]byte, error) {

	token, err := naverOauthConfig.Exchange(context.Background(), code)
	if err != nil {

		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	res, err := http.Get(oauthNaverUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll error %s\n ", err.Error())
	}

	return data, nil
}
