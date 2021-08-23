package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main Page",
	})
}

//MakeHandler is function to gather router.
func MakeHandler(dbname string) *gin.Engine {

	route := gin.Default()
	//Set up route groups and check session middleware
	route.GET("/", getIndex)

	//route.GET("/project/detail/iamport", Parentiamport)

	return route

}
