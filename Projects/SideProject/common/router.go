package common

import (
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

//ResArtistOfTheMonth is sturcture to contain response artist information.
type ResArtistOfTheMonth struct {
	Artist []Artist `json:"artist_list"`
}

//ArtistList is structure to get artist list information.
type Artist struct {
	ArtistID     int64  `json:"artist_id"`
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

func (p *project) GetIndexHandler(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main Page",
	})
}

func (p *project) GetProjectInfoHandler(c *gin.Context) {
	reqpod := &ReqProjectsOfTheDay{}
	err := c.ShouldBindJSON(reqpod)
	if err != nil {
		log.Println("[ERR] json err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	respod, err := p.db.ReadProjectList(reqpod)
	if err != nil {
		log.Println("[ERR] failed to read project list err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, respod)

}

func (p *project) GetArtistInfoHandler(c *gin.Context) {
	resaom, err := p.db.ReadArtistList()
	if err != nil {
		log.Println("[ERR] failed to read resaom  err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resaom)
}

func (p *project) GetUserInfoHandler(c *gin.Context) {

	userinfo := &User{}

	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user idd err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	resuser := p.db.ReadUserInfo(userinfo.UserID)
	c.JSON(http.StatusOK, resuser)
}

func (p *project) PutUserInfoHandler(c *gin.Context) {
	filefullpath, err := p.db.SaveJoinUserInfo(c)
	if err != nil {
		log.Println("[ERR] get image file err : ", err)
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
func (p *project) GetProfileFrameInfoHandler(c *gin.Context) {
	userinfo := &User{}

	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user idd err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	resprofileframeinfo, err := p.db.ReadProfileFrameInfo(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read profile frame information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofileframeinfo)
}
func (p *project) GetProfileProjectHandler(c *gin.Context) {
	userinfo := &User{}

	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user idd err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	resprofileprojectinfo, err := p.db.ReadProfileProjectInfo(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read profile project information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofileprojectinfo)
}

func (p *project) GetProfileSellHandler(c *gin.Context) {
	userinfo := &User{}

	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	resprofilesellinfo, err := p.db.ReadProfileSellInfo(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read profile sell history information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofilesellinfo)
}

func (p *project) GetProfileBuyHandler(c *gin.Context) {
	var userinfo = &User{}

	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	resprofilebuyinfo, err := p.db.ReadProfileBuyInfo(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read profile buy history information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofilebuyinfo)
}

func (p *project) GetProfileWithdrawHandler(c *gin.Context) {
	var userinfo = &User{}

	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	resprofilebuyinfo, err := p.db.ReadProfileWithdrawInfo(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read profile withdraw history information err :", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprofilebuyinfo)
}

//GetModificationUserInfoHandler is function to get user profile infomation.
func (p *project) GetModificationUserInfoHandler(c *gin.Context) {

	var userinfo = &User{}
	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	resmodificationuserinfo, err := p.db.ReadModificationUserInfo(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read modification user info err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, resmodificationuserinfo)
}

//PutModificationUserInfoHandler is fuction to edit user profile information.
func (p *project) PutModificationUserInfoHandler(c *gin.Context) {
	var reqmodinfo = &ReqModificationUserInfo{}
	err := c.ShouldBindJSON(reqmodinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	err = p.db.UpdateModificationUserInfo(reqmodinfo)
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
	var artistprofileinfo = &ArtistProfile{}
	err := c.ShouldBindJSON(artistprofileinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	resartistinfo, err := p.db.ReadProfileArtistInfo(artistprofileinfo.ArtistID)
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
	var userinfo = &User{}
	err := c.ShouldBindJSON(userinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	respersonalinfomation, err := p.db.ReadPersonalInformation(userinfo.UserID)
	if err != nil {
		log.Println("[ERR] failed to read personal information err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, respersonalinfomation)
}

func (p *project) PutPersonalInformationHandler(c *gin.Context) {
	var reqpersonalinformation = &ReqPersonalInformation{}
	err := c.ShouldBindJSON(reqpersonalinformation)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	err = p.db.UpdatePersonalInformation(reqpersonalinformation)
	if err != nil {
		log.Println("[ERR] failed to update personal information err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (p *project) GetProjectDetailProjectInformationHandler(c *gin.Context) {
	var reqprojectdetailinfo = &ReqProjectDetailInfo{}
	err := c.ShouldBindJSON(reqprojectdetailinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	resprojectdetail, err := p.db.ReadProjectDetailArtistProjectInfo(reqprojectdetailinfo)
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
		profileuser.GET("/frame", CheckSessionValidity, p.GetProfileFrameInfoHandler)
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
		//------ {TODO : project and frame split
		profileartist.GET("/information", CheckSessionValidity, p.GetProfileArtistInfoHandler)
		//------ TODO : project and frame split }
	}

	personal := router.Group("/personal")
	{
		personal.GET("/index", p.GetPersonalIndexHandler)
		personal.GET("/information", p.GetPersonalInformationHandler)
		personal.PUT("/information", p.PutPersonalInformationHandler)
	}

	return router

}
