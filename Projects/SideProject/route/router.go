package route

import (
	"net/http"
	"sideproject/route/artist"
	"sideproject/route/auth"
	"sideproject/route/database"
	"sideproject/route/iamport"
	"sideproject/route/profile"
	"sideproject/route/project"
	"sideproject/route/user"

	"github.com/gin-gonic/gin"

	"github.com/unrolled/render"
)

var rd *render.Render = render.New()

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
	i := &iamport.Iamport{DB: d, APIKey: iamport.APIKey, APISecret: iamport.APISecret}

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
