package artist

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Artist structure is used to seperate that these are Artist's functions .
type Artist struct {
	DB *sql.DB
}

//ArtistList is structure to get artist list information.
type JSArt struct {
	ArtistID     int64  `json:"artist_id"`
	NickName     string `json:"nickname"`
	Introduction string `json:"introduction"`
	ImageLink    string `json:"image_link"`
	Rank         string `json:"rank"`
}

//ResJSArtOfTheMonth is sturcture to contain response artist information.
type ResJSArtOfTheMonth struct {
	Art []JSArt `json:"artist_list"`
}

func (a *Artist) Routes(route *gin.RouterGroup) {
	route.GET("/", a.getArtistInfo)
}

func (a *Artist) getArtistInfo(c *gin.Context) {
	resaom, err := a.ReadArtistList()
	if err != nil {
		log.Println("[ERR] failed to read resaom  err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resaom)
}

func (a *Artist) ReadArtistList() (*ResJSArtOfTheMonth, error) {

	var resaom *ResJSArtOfTheMonth
	resaom = &ResJSArtOfTheMonth{}

	stmt, err := a.DB.Prepare(`SELECT 
								u.id
								u.nickname,
								u.introduction,
								u.image_link,
								RANK() OVER(ORDER BY u_r.score DESC)
								FROM user AS u INNER JOIN user_rank AS u_r ON u.id = u_r.user_id 
								WHERE DATE_FORMAT(u_r.created_at,"%Y-%m") = DATE_FORMAT(NOW(),"%Y-%m") LIMIT 5;`)

	if err != nil {
		log.Println("[ERR] prepare statement err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Println("[ERR] query err : ", err)
		return nil, err
	}

	defer rows.Close()

	var jsart JSArt
	for rows.Next() {
		rows.Scan(&jsart.ArtistID, &jsart.NickName, &jsart.Introduction, &jsart.ImageLink, &jsart.Rank)
		resaom.Art = append(resaom.Art, jsart)
	}

	return resaom, nil
}
