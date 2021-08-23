package route

import (
	"io/ioutil"
	"log"
	"net/http"
	"sideproject/route/artist"
	"sideproject/route/auth"
	"sideproject/route/database"
	"sideproject/route/iamport"
	"sideproject/route/profile"
	"sideproject/route/project"
	"sideproject/route/user"
	"time"

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

func getIndex(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main Page",
	})
}

//MakeHandler is function to gather router.
func MakeHandler(dbname string) *gin.Engine {

	d := database.MakeDBHandler(dbname)

	u := &user.Usr{DB: d}
	pf := &profile.Profile{DB: d}
	p := &project.Pjt{DB: d}
	a := &artist.Artist{DB: d}
	au := &auth.Auth{DB: d}
	i := &iamport.Iamport{DB: d}

	route := gin.Default()
	//Set up route groups and check session middleware
	route.GET("/", getIndex)

	grusr := route.Group("/user")
	grusr.Use(au.CheckSessionValidity)

	grpf := route.Group("/profile")
	grpf.Use(au.CheckSessionValidity)

	grp := route.Group("/project")
	grp.Use(au.CheckSessionValidity)

	gra := route.Group("/artist")
	gra.Use(au.CheckSessionValidity)

	gi := route.Group("/iamport")
	gi.Use(au.CheckSessionValidity)

	u.Routes(grusr)
	pf.Routes(grpf)
	p.Routes(grp)
	a.Routes(gra)
	i.Routes(gi)

	return route

}

func Parentiamport(c *gin.Context) {

	body, err := Iamport()
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, body)

}

func Iamport() ([]byte, error) {

	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, "https://admin.iamport.kr/users/getToken", nil)
	if err != nil {
		log.Println("[ERR] new request err : ", err)
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
		return nil, err
	}

	body, _ := ioutil.ReadAll(res.Body)

	return body, nil
}
