package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/unrolled/render"

	artist "sideproject.com/artist"
	auth "sideproject.com/auth"
	database "sideproject.com/database"
	profile "sideproject.com/profile"
	project "sideproject.com/project"
	user "sideproject.com/user"
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
func MakeHandler(route *gin.Engine, dbname string) *gin.Engine {

	d := database.MakeDBHandler(dbname)

	u := &user.Usr{DB: d}
	pf := &profile.Profile{DB: d}
	p := &project.Pjt{DB: d}
	a := &artist.Artist{DB: d}
	au := &auth.Auth{DB: d}

	//Set up route groups and check session middleware
	grusr := route.Group("/user")
	grusr.Use(au.CheckSessionValidity)

	grpf := route.Group("/profile")
	grpf.Use(au.CheckSessionValidity)

	grp := route.Group("/project")
	grp.Use(au.CheckSessionValidity)

	gra := route.Group("/artist")
	gra.Use(au.CheckSessionValidity)

	route.GET("/", getIndex)

	u.Routes(grusr)
	pf.Routes(grpf)
	p.Routes(grp)
	a.Routes(route)

	return route

}
