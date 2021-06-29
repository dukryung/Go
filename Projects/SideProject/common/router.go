package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/unrolled/render"
)

var rd *render.Render = render.New()

const htmlIndex = `<html><body>
Logged in with <a href="/auth/google/signup">google</a>
</br>
Logged in with <a href="/auth/facebook/signup">facebook</a>
</br>
Logged in with <a href="/auth/kakao/signup">kakao</a>
</br>
Logged in with <a href="/auth/naver/signup">naver</a>
</br>
</body></html>
`

type project struct {
	db DBHandler
}

//ReqProjectsOfTheDay is structure to contain request information.
type ReqProjectsOfTheDay struct {
	DemandDate   string `json:"demand_date"`
	DemandPeriod string `json:"demand_period"`
}

//ResProjectsOfTheDay is structure to contain response information.
type ResProjectsOfTheDay struct {
	Date           string        `json:"date"`
	Project        []ProjectList `json:"project"`
	Total          string        `json:"total"`
	Period         string        `json:"period"`
	RankLastNumber string        `json:"rank_last_number"`
}

//ProjectList is structure to get project list information.
type ProjectList struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	CategoryCode string `json:"category_code"`
	Description  string `json:"desc"`
	ImageLink    string `json:"image_link"`
	CreatedAt    string `json:"created_at"`
	SellCount    string `json:"sell_cnt"`
	UserNickName string `json:"user_nickname"`
	CommentCount string `json:"comment_count"`
	UpvoteCount  string `json:"upvote_count"`
	Price        string `json:"price"`
	Beta         string `json:"beta"`
	Rank         string `json:"rank"`
}

//ResArtistOfTheMonth is sturcture to contain response artist information.
type ResArtistOfTheMonth struct {
	Artist []ArtistList `json:"artist"`
}

//ArtistList is structure to get artist list information.
type ArtistList struct {
	NickName     string `json:"nickname"`
	Introduction string `json:"introduction"`
	ImageLink    string `json:"image_link"`
	Rank         string `json:"rank"`
}

//ResUserInfo is structure to contain user information in index page.
type ResUserInfo struct {
	ID       string `json:"id"`
	NickName string `json:"nickname"`
}

type ResEmailInfo struct {
	Email string `json:"email"`
}

type ReqJoinInfo struct {
	UserInfo    User    `json:"user"`
	AccountInfo Account `json:"account"`
}
type User struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Nickname            string `json:"nickname"`
	Email               string `json:"email"`
	Introduction        string `json:"introduction"`
	AgreeEmailMarketing bool   `json:"agree_email_marketing"`
}

type Account struct {
	UserID      string `json:"user_id"`
	Bank        int    `json:"bank"`
	Account     string `json:"account"`
	AgreePolicy bool   `json:"agree_policy"`
}

func (p *project) GetProjectInfoHandler(c *gin.Context) {
	var reqpod ReqProjectsOfTheDay

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("[ERR] readAll err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
	}

	err = json.Unmarshal(data, &reqpod)
	if err != nil {
		log.Println("[ERR] json unmarshal err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
	}

	log.Println("[LOG] reqpod information : ", reqpod)

	respod := p.db.ReadProjectList(&reqpod)
	if respod == nil {
		log.Println("[LOG] empty project list")
		c.JSON(http.StatusInternalServerError, nil)
	}

	c.JSON(http.StatusOK, respod)

}

func (p *project) GetArtistInfoHandler(c *gin.Context) {
	resaom := p.db.ReadArtistList()
	if resaom == nil {
		log.Println("[LOG] resaom : ", resaom)
		log.Println("[LOG] empty artist list")
	}
	c.JSON(http.StatusOK, resaom)
}

func (p *project) GetUserInfoHandler(c *gin.Context) {

	session, _ := store.Get(c.Request, "session")
	val := session.Values["id"]
	sessionid := val.(string)

	resuser := p.db.ReadUserInfo(sessionid)
	c.JSON(http.StatusOK, resuser)
}

func (p *project) GetIndexHandler(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main Page",
	})
}

func (p *project) PutUserInfoHandler(c *gin.Context) {

	var img Image

	filefullpath, reqjoininfo, err := img.SaveImageFiles(c)
	if err != nil {
		log.Println("[ERR] get image file err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	log.Println("[LOG] img link array : ", filefullpath)

	err = p.db.UpdateUserInfo(filefullpath, reqjoininfo)
	if err != nil {
		log.Println("[ERR] update user info err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filedirinfo": filefullpath,
	})
}

func (p *project) GetProfileHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "profileuser.html", gin.H{
		"title": "Profile User Page",
	})
}
func (p *project) GetProfileFramInfoHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	val := session.Values["id"]
	sessionid := val.(string)

	resprofileframeinfo, err := p.db.ReadProfileFrameInfo(sessionid)
	if err != nil {
		log.Println("[ERR] failed to read profile frame information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofileframeinfo)
}
func (p *project) GetProfileProjectHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	val := session.Values["id"]
	sessionid := val.(string)

	resprofileprojectinfo, err := p.db.ReadProfileProjectInfo(sessionid)
	if err != nil {
		log.Println("[ERR] failed to read profile project information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofileprojectinfo)
}

func (p *project) GetProfileSellHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	val := session.Values["id"]
	sessionid := val.(string)

	resprofilesellinfo, err := p.db.ReadProfileSellInfo(sessionid)
	if err != nil {
		log.Println("[ERR] failed to read profile sell history information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofilesellinfo)
}

func (p *project) GetProfileBuyHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	val := session.Values["id"]
	sessionid := val.(string)

	resprofilebuyinfo, err := p.db.ReadProfileBuyInfo(sessionid)
	if err != nil {
		log.Println("[ERR] failed to read profile buy history information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofilebuyinfo)
}

func (p *project) GetProfileWithdrawHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	val := session.Values["id"]
	sessionid := val.(string)

	resprofilebuyinfo, err := p.db.ReadProfileWithdrawInfo(sessionid)
	if err != nil {
		log.Println("[ERR] failed to read profile withdraw history information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofilebuyinfo)
}

//GetModificationUserInfoHandler is function to get user profile infomation.
func (p *project) GetModificationUserInfoHandler(c *gin.Context) {

	session, _ := store.Get(c.Request, "session")
	val := session.Values["id"]
	sessionid := val.(string)

	resmodificationuserinfo, err := p.db.ReadModificationUserInfo(sessionid)
	if err != nil {
		log.Println("[ERR] failed to read modification user info err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, resmodificationuserinfo)
}

//PutProfileEditHandler is fuction to edit user profile information.
func (p *project) PutModificationUserInfoHandler(c *gin.Context) {
	var reqjoininfo *ReqJoinInfo

	session, _ := store.Get(c.Request, "session")
	val := session.Values["id"]
	sessionid := val.(string)

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("[ERR] ioutil readall err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	err = json.Unmarshal(data, reqjoininfo)
	if err != nil {
		log.Println("[ERR] json unmarshal err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	err = p.db.UpdateModificationUserInfo(sessionid, reqjoininfo)
	if err != nil {
		log.Println("[ERR] update user info ")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, nil)

}
func (p *project) GetProfileArtistHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "profileartist.html", gin.H{
		"title": "profile artist  Page",
	})
}

func (p *project) GetProfileArtistInfoHandler(c *gin.Context) {
	resartistinfo, err := p.db.ReadProfileArtistInfo(c)
	if err != nil {
		log.Println("[ERR] failed to read artist information err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resartistinfo)
}

func (p *project) GetPersonalIndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "personalinformation.html", gin.H{
		"title": "Personal Information Page",
	})
}
func (p *project) GetPersonalInformationHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	val := session.Values["id"]
	sessionid := val.(string)

	respersonalinfomation, err := p.db.ReadPersonalInformation(sessionid)
	if err != nil {
		log.Println("[ERR] failed to read personal information err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, respersonalinfomation)
}

func (p *project) PutPersonalInformationHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	val := session.Values["id"]
	sessionid := val.(string)

	err := p.db.UpdatePersonalInformation(c, sessionid)
	if err != nil {
		log.Println("[ERR] failed to update personal information err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (p *project) GetProjectDetailProjectInformationHandler(c *gin.Context) {

	resprojectdetail, err := p.db.ReadProjectDetailArtistProjectInfo(c)
	if err != nil {
		log.Println("[ERR] failed to read project detail information err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprojectdetail)

}

func (p *project) GetProjectDetailProjectImagesHandler(c *gin.Context) {

	resprojectdetailimagesinfo, err := p.db.ReadProjectDetailArtistProjectImagesInfo(c)
	if err != nil {
		log.Println("[ERR] failed to read project detail image links err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprojectdetailimagesinfo)
}

func (p *project) GetProjectDetailCommentHandler(c *gin.Context) {
	resprojectdetailimagesinfo, err := p.db.ReadProjectDetailCommentInfo(c)
	if err != nil {
		log.Println("[ERR] failed to read project detail image links err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, resprojectdetailimagesinfo)
}

func (p *project) PostProjectUploadHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	val := session.Values["id"]
	sessionid := val.(string)

	userid, err := p.db.ReadUserID(sessionid)
	if err != nil {
		log.Println("[ERR] failed to read userid err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	err = p.db.CreateProjectInfo(c, userid)
	if err != nil {
		log.Println("[ERR] failed to upload project err : ", err)
		c.JSON(http.StatusOK, nil)
		return
	}

	c.JSON(http.StatusOK, nil)

}

func (p *project) PutProjectUploadHandler(c *gin.Context) {

}

//MakeHandler is function to gather router.
func MakeHandler(databasename string) *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("../public/*")
	pdb := &ProjectDB{}
	p := &project{db: MakeDBHandler(databasename, pdb)}

	router.GET("/", p.GetIndexHandler)

	project := router.Group("/project")
	{
		project.GET("/information", p.GetProjectInfoHandler)
		project.GET("/information/detail/project", p.GetProjectDetailProjectInformationHandler)
		project.GET("/informaiton/detail/image", p.GetProjectDetailProjectImagesHandler)
		project.GET("/informaiton/detail/comment", p.GetProjectDetailCommentHandler)
		project.POST("/information/upload")
		project.PUT("/information/upload")
	}

	router.GET("/artist", p.GetArtistInfoHandler)

	router.GET("/user", CheckSessionValidity, p.GetUserInfoHandler)
	router.PUT("/user", CheckSessionValidity, p.PutUserInfoHandler)

	authorization := router.Group("/auth")
	{
		authorization.GET("/google/signup", p.googleLoginHandler)
		authorization.GET("/google/callback", p.googleAuthCallback)
		authorization.GET("/facebook/signup", p.facebookLoginHandler)
		authorization.GET("/facebook/callback", p.facebookAuthCallback)
		authorization.GET("/kakao/signup", p.kakaoLoginHandler)
		authorization.GET("/kakao/callback", p.kakaoAuthCallback)
		authorization.GET("/naver/signup", p.naverLoginHndler)
		authorization.GET("/naver/callback", p.naverAuthCallback)
	}

	profileuser := router.Group("/profileuser")
	{
		profileuser.GET("/index", CheckSessionValidity, p.GetProfileHandler)
		profileuser.GET("/frame", CheckSessionValidity, p.GetProfileFramInfoHandler)
		profileuser.GET("/project", CheckSessionValidity, p.GetProfileProjectHandler)
		profileuser.GET("/sell", CheckSessionValidity, p.GetProfileSellHandler)
		profileuser.GET("/buy", CheckSessionValidity, p.GetProfileBuyHandler)
		profileuser.GET("/withdraw", CheckSessionValidity, p.GetProfileWithdrawHandler)
		profileuser.GET("/modification", CheckSessionValidity, p.GetModificationUserInfoHandler)
		profileuser.PUT("/modification", CheckSessionValidity, p.PutModificationUserInfoHandler)

	}

	profileartist := router.Group("/profileartist")
	{
		profileartist.GET("/index", CheckSessionValidity, p.GetProfileArtistHandler)
		profileartist.GET("/information", CheckSessionValidity, p.GetProfileArtistInfoHandler)
	}

	personal := router.Group("/personal")
	{
		personal.GET("/index", p.GetPersonalIndexHandler)
		personal.GET("/information", p.GetPersonalInformationHandler)
		personal.PUT("/information", p.PutPersonalInformationHandler)
	}

	return router

}