package project

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"sideproject/route/gcloud"
	"strconv"
	"time"

	"google.golang.org/appengine"

	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

type Pjt struct {
	DB *sql.DB
}
type ReqProjectDetailInfo struct {
	UserID    int `json:"user_id"`
	ProjectID int `json:"project_id"`
}

type ResProjectDetailInfo struct {
	ProjectDetailInfo           ProjectDetail                `json:"projectdetail"`
	ProjectDetailImageLinksInfo []ProjectDetailImageLinkInfo `json:"image_links"`
	ProjectCommentInfo          []ProjectDetailCommentInfo   `json:"comments"`
}

type ProjectDetail struct {
	ArtistID       int    `json:"aritst_id"`
	ArtistNickname string `json:"aritst_nickname"`
	Title          string `json:"title"`
	Desc           string `json:"desc"`
	VideoLink      string `json:"video_link"`
	UpvoteCount    int    `json:"upvote_count"`
	CreatedAt      string `json:"created_at"`
	Rank           int    `json:"rank"`
	BetaLink       string `json:"beta_link"`
	UpvoteStatus   bool   `json:"upvote_status"`
	BuyStatus      bool   `json:"buy_status"`
}

type ProjectDetailImageLinkInfo struct {
	ImageLink string `json:"image_link"`
}

type ProjectDetailCommentInfo struct {
	ID             int                      `json:"id"`
	ArtistNickName string                   `json:"artist_nickname"`
	Text           string                   `json:"text"`
	Time           string                   `json:"time"`
	Replies        []ProjectDetailReplyInfo `json:"replies"`
}
type ProjectDetailReplyInfo struct {
	ID             int    `json:"id"`
	ArtistNickName string `json:"artist_nickname"`
	Text           string `json:"text"`
	Time           string `json:"time"`
}

type ReqUploadProjectInfo struct {
	ID          int64  `json:"id"`
	ProjectID   int64  `json:"projectid"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Desc        string `json:"desc"`
	Price       int    `json:"price"`
	YoutubeLink string `json:"youtube"`
}

//ReqProjectsOfTheDay is structure to contain request informationå.
type ReqProjectsOfTheDay struct {
	DemandDate   time.Time `json:"demand_date" binding:"required"`
	DemandPeriod int       `json:"demand_period" binding:"required"`
}

//ResProjectsOfTheDay is structure to contain response information.
type ResProjectsOfTheDay struct {
	Date           time.Time     `json:"date"`
	Project        []ProjectList `json:"project_list"`
	Total          string        `json:"total"`
	Period         int           `json:"period"`
	RankLastNumber string        `json:"rank_last_number"`
}

//ProjectList is structure to get project list information.
type ProjectList struct {
	ID            int64  `json:"project_id"`
	Title         string `json:"title"`
	CategoryCode  string `json:"category_id"`
	Description   string `json:"desc"`
	ImageLink     string `json:"image_link"`
	CreatedAt     string `json:"created_at"`
	SellCount     string `json:"sell_cnt"`
	AristNickName string `json:"artist_nickname"`
	CommentCount  string `json:"comment_count"`
	UpvoteCount   string `json:"upvote_count"`
	Price         string `json:"price"`
	Beta          string `json:"beta"`
	Rank          string `json:"rank"`
}

func (p *Pjt) Routes(route *gin.RouterGroup) {

	route.GET("/information", p.getProject)
	route.GET("/information/detail/project", p.getProjectDetailProjectInformation)
	route.GET("/informaiton/detail/image", p.getProjectDetailProjectImages)
	route.GET("/informaiton/detail/comment", p.getProjectDetailComment)
	route.POST("/information/upload", p.postProjectUpload)
	route.PUT("/information/upload", p.putProjectUpload)

}
func (p *Pjt) putProjectUpload(c *gin.Context) {
	err := p.UpdateProjectInfo(c)
	if err != nil {
		log.Println("[ERR] failed to update project err : ", err)
		c.JSON(http.StatusOK, nil)
		return
	}

	c.JSON(http.StatusOK, nil)

}

func (p *Pjt) postProjectUpload(c *gin.Context) {

	err := p.CreateProjectInfo(c)
	if err != nil {
		log.Println("[ERR] failed to upload project err : ", err)
		c.JSON(http.StatusOK, nil)
		return
	}

	c.JSON(http.StatusOK, nil)

}

func (p *Pjt) getProjectDetailProjectInformation(c *gin.Context) {
	var reqprojectdetailinfo = &ReqProjectDetailInfo{}
	err := c.ShouldBindJSON(reqprojectdetailinfo)
	if err != nil {
		log.Println("[ERR] failed to extract user id err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	resprojectdetail, err := p.ReadProjectDetailArtistProjectInfo(reqprojectdetailinfo)
	if err != nil {
		log.Println("[ERR] failed to read project detail information err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprojectdetail)

}

func (p *Pjt) getProjectDetailProjectImages(c *gin.Context) {

	resprojectdetailimagesinfo, err := p.ReadProjectDetailArtistProjectImagesInfo(c)
	if err != nil {
		log.Println("[ERR] failed to read project detail image links err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, resprojectdetailimagesinfo)
}

func (p *Pjt) getProjectDetailComment(c *gin.Context) {
	resprojectdetailimagesinfo, err := p.ReadProjectDetailCommentInfo(c)
	if err != nil {
		log.Println("[ERR] failed to read project detail image links err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, resprojectdetailimagesinfo)
}

func (p *Pjt) getProject(c *gin.Context) {
	reqpod := &ReqProjectsOfTheDay{}
	err := c.ShouldBindJSON(reqpod)
	if err != nil {
		log.Println("[ERR] json err : ", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	respod, err := p.ReadProjectList(reqpod)
	if err != nil {
		log.Println("[ERR] failed to read project list err : ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, respod)

}

func (p *Pjt) ReadProjectList(reqpod *ReqProjectsOfTheDay) (*ResProjectsOfTheDay, error) {

	var respod *ResProjectsOfTheDay
	respod = &ResProjectsOfTheDay{Date: reqpod.DemandDate}

	stmt, err := p.DB.Prepare(`SELECT 
						p.id, 
						p.title, 
						p.category_id,
				 		p.description,
						DATE_FORMAT(p.created_at,"%Y-%m-%d"),
						p.sell_count, 
						u.nickname, 
						p.comment_count,
						p.total_upvote_count,
						p.price,
						p.beta,
						RANK() OVER(ORDER BY p_r.score DESC)   
				  FROM project AS p INNER JOIN project_rank AS p_r ON p.id = p_r.project_id 
				  					INNER JOIN user AS u ON p.user_id = u.id
				  WHERE p_r.created_at BETWEEN DATE_FORMAT(DATE_SUB( ? ,INTERVAL ? DAY),"%Y-%m-%d") AND DATE_FORMAT(DATE_ADD( ? ,INTERVAL ? DAY), "%Y-%m-%d")
				  LIMIT 10;`)

	if err != nil {
		log.Println("[ERR] stmt err : ", err)
		return nil, err
	}
	defer stmt.Close()

	var rows *sql.Rows

	if reqpod.DemandPeriod == 1 {
		rows, err = stmt.Query(reqpod.DemandDate, reqpod.DemandPeriod-1, reqpod.DemandDate, reqpod.DemandPeriod)
		if err != nil {
			log.Println("[ERR] rows err : ", err)
			return nil, err
		}
	} else if reqpod.DemandPeriod == 7 || reqpod.DemandPeriod == 30 {
		rows, err = stmt.Query(reqpod.DemandDate, reqpod.DemandPeriod-1, reqpod.DemandDate, 1)
		if err != nil {
			log.Println("[ERR] rows err : ", err)
			return nil, err
		}
	} else {
		return nil, errors.New("[ERR] wrong period")
	}

	defer rows.Close()

	var project ProjectList
	respod.Date = reqpod.DemandDate
	respod.Period = reqpod.DemandPeriod
	respod.Total = "0"
	respod.RankLastNumber = "0"

	//var projectid, ranking uint64
	var projectid int64
	var ranking, title, categorycode, desc, createdat, sellcount, artistnickname, commentcount, totalupvotecount, price, beta string
	for rows.Next() {
		err := rows.Scan(&projectid, &title, &categorycode, &desc, &createdat, &sellcount, &artistnickname, &commentcount, &totalupvotecount, &price, &beta, &ranking)
		if err != nil {
			log.Println("[ERR] scan err : ", err)
			return nil, err
		}

		project.ID = projectid
		project.Title = title
		project.CategoryCode = categorycode
		project.Description = desc
		project.CreatedAt = createdat
		project.SellCount = sellcount
		project.AristNickName = artistnickname
		project.CommentCount = commentcount
		project.UpvoteCount = totalupvotecount
		project.Price = price
		project.Beta = beta
		project.Rank = ranking

		ranktypeuint64, _ := strconv.ParseUint(ranking, 10, 32)
		respod.RankLastNumber = strconv.FormatUint(ranktypeuint64+1, 10)
		respod.Project = append(respod.Project, project)
	}

	var link string
	for i := range respod.Project {
		stmt, err = p.DB.Prepare(`SELECT link FROM image WHERE id = ? LIMIT 1;`)
		rows, err = stmt.Query(respod.Project[i].ID)
		if err != nil {
			log.Println("[ERR] stmt query err : ", err)
			return nil, err
		}

		for rows.Next() {
			err := rows.Scan(&link)
			if err != nil {
				log.Println("[ERR] rows scan err : ", err)
				return nil, err
			}
			respod.Project[i].ImageLink = link
		}
	}

	stmt, err = p.DB.Prepare(`SELECT COUNT(*) 
				  FROM project 
				  WHERE created_at 
				  BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL ? DAY),"%Y-%m-%d") AND DATE_FORMAT(DATE_ADD(?, INTERVAL ? DAY),"%Y-%m-%d")`)

	if reqpod.DemandPeriod == 1 {
		rows, err = stmt.Query(reqpod.DemandDate, reqpod.DemandPeriod-1, reqpod.DemandDate, reqpod.DemandPeriod)
		if err != nil {
			log.Println("[ERR] rows err : ", err)
			return nil, err
		}
	} else if reqpod.DemandPeriod == 7 || reqpod.DemandPeriod == 30 {
		rows, err = stmt.Query(reqpod.DemandDate, reqpod.DemandPeriod-1, reqpod.DemandDate, 1)
		if err != nil {
			log.Println("[ERR] rows err : ", err)
			return nil, err
		}
	} else {
		return nil, errors.New("[ERR] wrong period")
	}

	var projectcnt string
	for rows.Next() {
		err := rows.Scan(&projectcnt)
		if err != nil {
			return nil, err
		}
		respod.Total = projectcnt
	}

	return respod, nil
}

func (p *Pjt) ReadProjectDetailArtistProjectInfo(reqprojectdetailinfo *ReqProjectDetailInfo) (*ResProjectDetailInfo, error) {

	var resprojectdetailinfo *ResProjectDetailInfo

	// -------{TODO : delete rank column and test query  in database's table.}
	stmt, err := p.DB.Prepare(`SELECT 
				  u.id,
				  u.nickname, 
				  p.title,
				  p.desc, 
				  p.video_link,
				  p.total_upvote_count,
				  p.created_at,
				  pr.rank, 
				  p.beta_link,
				  EXISTS(us.project_id) as SUCCESS,
				  EXISTS(bh.project_id) as SUCCESS
				  FROM user AS u INNER JOIN project AS p ON p.user_id = u.id
				  INNER JOIN project_rank AS pr ON pr.project_id = p.id
				  INNER JOIN upvote_status AS us ON us.user_id = u.id
				  INNER JOIN buy_history AS bh ON bh.user_id = u.id 
				  WHERE u.id = ? AND p.id = ?`)

	if err != nil {
		log.Println("[ERR] prepare statement err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(reqprojectdetailinfo.UserID, reqprojectdetailinfo.ProjectID)
	if err != nil {
		log.Println("[ERR] statement query err : ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(resprojectdetailinfo.ProjectDetailInfo.ArtistID, resprojectdetailinfo.ProjectDetailInfo.ArtistNickname, resprojectdetailinfo.ProjectDetailInfo.Title, resprojectdetailinfo.ProjectDetailInfo.Desc, resprojectdetailinfo.ProjectDetailInfo.VideoLink, resprojectdetailinfo.ProjectDetailInfo.UpvoteCount, resprojectdetailinfo.ProjectDetailInfo.CreatedAt, resprojectdetailinfo.ProjectDetailInfo.Rank, resprojectdetailinfo.ProjectDetailInfo.BetaLink, resprojectdetailinfo.ProjectDetailInfo.UpvoteStatus, resprojectdetailinfo.ProjectDetailInfo.BuyStatus)
		if err != nil {
			log.Println("[ERR] rows scan err : ", err)
			return nil, err
		}
	}

	return resprojectdetailinfo, nil
}

func (p *Pjt) ReadProjectDetailArtistProjectImagesInfo(c *gin.Context) (*ResProjectDetailInfo, error) {
	var reqprojectdtailinfo *ReqProjectDetailInfo
	var projectdetailimagelinkinfo ProjectDetailImageLinkInfo
	var resprojectdetailinfo *ResProjectDetailInfo

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("[ERR] read all err : ", err)
		return nil, err
	}

	err = json.Unmarshal(data, reqprojectdtailinfo)
	if err != nil {
		log.Println("[ERR] json unmarshal err : ", err)
		return nil, err
	}

	stmt, err := p.DB.Prepare(`SELECT 
				  link
				  FROM image 
				  WHERE project_id = ?
				  `)

	if err != nil {
		log.Println("[ERR] prepare statement err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(reqprojectdtailinfo.ProjectID)
	if err != nil {
		log.Println("[ERR] statement err :", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(projectdetailimagelinkinfo.ImageLink)
		if err != nil {
			log.Println("[ERR] rows scan err :", err)
			return nil, err
		}
		resprojectdetailinfo.ProjectDetailImageLinksInfo = append(resprojectdetailinfo.ProjectDetailImageLinksInfo, projectdetailimagelinkinfo)
	}

	return resprojectdetailinfo, nil

}

func (p *Pjt) ReadProjectDetailCommentInfo(c *gin.Context) (*ResProjectDetailInfo, error) {
	var reqprojectdtailinfo *ReqProjectDetailInfo
	var resprojectdetailinfo *ResProjectDetailInfo
	var projectdetailcommentinfo ProjectDetailCommentInfo
	var projectdetailreplyinfo ProjectDetailReplyInfo

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("[ERR] read all err : ", err)
		return nil, err
	}

	err = json.Unmarshal(data, reqprojectdtailinfo)
	if err != nil {
		log.Println("[ERR] json unmarshal err : ", err)
		return nil, err
	}

	stmt, err := p.DB.Prepare(`SELECT
				  c.artist_id,
				  c.artist_nickname,
				  c.text,
				  c.created_at,
				  r.artust_id,
				  r.artist_nickname,
				  r.text,
				  r.created_at
				  FROM comment AS c 
				  INNER JOIN reply AS r ON r.comment_id = c.id
				  WHERE c.project_id = ?`)
	if err != nil {
		log.Println("[ERR] prepare statement err : ", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(reqprojectdtailinfo.ProjectID)
	if err != nil {
		log.Println("[ERR] statement err : ", err)
		return nil, err
	}

	defer rows.Close()

	var compareartistid int
	for rows.Next() {
		rows.Scan(projectdetailcommentinfo.ID, projectdetailcommentinfo.ArtistNickName, projectdetailcommentinfo.Text, projectdetailcommentinfo.Time, projectdetailreplyinfo.ID, projectdetailreplyinfo.ArtistNickName, projectdetailreplyinfo.Text, projectdetailreplyinfo.Time)

		if compareartistid != projectdetailcommentinfo.ID {
			resprojectdetailinfo.ProjectCommentInfo = append(resprojectdetailinfo.ProjectCommentInfo, projectdetailcommentinfo)
		}

		for i := range resprojectdetailinfo.ProjectCommentInfo {
			if len(resprojectdetailinfo.ProjectCommentInfo) == i+1 {
				resprojectdetailinfo.ProjectCommentInfo[i].Replies = append(resprojectdetailinfo.ProjectCommentInfo[i].Replies, projectdetailreplyinfo)
			}
		}
	}

	return resprojectdetailinfo, nil
}

func (p *Pjt) CreateProjectInfo(c *gin.Context) error {
	var isbeta = false
	var videoobject gcloud.Video
	var projectfile gcloud.ProjectFile
	var image gcloud.Image
	var lastid int64
	var projectimagelinks []string
	var videolink, betalink, originlink string

	requploadprojectinfo, videofd, imagefdarr, betafd, originfd, err := GetProjectInfo(c)
	if err != nil {
		log.Println("[ERR] failed to get project informaiton err : ", err)
		return err
	}

	ctx := appengine.NewContext(c.Request)

	tx, err := p.DB.Begin()
	if err != nil {
		log.Println("[ERR] begin err : ", err)
		return err
	}

	defer tx.Rollback()

	{
		stmt, err := tx.Prepare(`INSERT INTO project (user_id, category_id, title, description, price, video_link, beta, beta_link, created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,NOW(),NOW())`)
		if err != nil {
			log.Println("[ERR] prepare statement err : ", err)
			return err
		}
		defer stmt.Close()

		result, err := stmt.Exec(requploadprojectinfo.ID, requploadprojectinfo.Category, requploadprojectinfo.Title, requploadprojectinfo.Desc, requploadprojectinfo.Price, requploadprojectinfo.YoutubeLink, isbeta, betalink)
		if err != nil {
			log.Println("[ERR] stmt exec err : ", err)
			return err
		}

		rowcnt, err := result.RowsAffected()
		if err != nil {
			log.Println("[ERR] rows affected err : ", err)
			return err
		}

		if rowcnt == 0 {
			log.Println("[LOG] None affected rows")
		}

		lastid, err = result.LastInsertId()
		if err != nil {
			log.Println("[ERR] last insert id err : ", err)
			return err
		}
	}

	{
		originlink, err = projectfile.SaveOriginFile(ctx, requploadprojectinfo.ID, lastid, originfd)
		if err != nil {
			log.Println("[ERR] failed to save origin file err : ", err)
			return err
		}

		stmt, err := tx.Prepare(`UPDATE project SET origin_link = ? WHERE id = ?`)
		if err != nil {
			log.Println("[ERR] prepare statement err : ", err)
			return err
		}
		defer stmt.Close()

		result, err := stmt.Exec(originlink, lastid)
		if err != nil {
			log.Println("[ERR] stmt exec err : ", err)
			return err
		}

		rowcnt, err := result.RowsAffected()
		if err != nil {
			log.Println("[ERR] rows affected err : ", err)
			return err
		}

		if rowcnt == 0 {
			log.Println("[LOG] None affected rows")
		}

	}

	{
		if videofd != nil {
			videolink, err = videoobject.SaveVideoFile(ctx, requploadprojectinfo.ID, lastid, videofd)
			if err != nil {
				log.Println("[ERR] failed to save video file  err : ", err)
				err = gcloud.DeleteOriginFile(ctx, requploadprojectinfo.ID, lastid, p.DB)
				if err != nil {
					log.Println("[ERR] failed to delete origin file err : ", err)
					return err
				}
				return err
			}

		}
		stmt, err := tx.Prepare(`UPDATE project SET video_link=? WHERE id = ?`)
		if err != nil {
			log.Println("[ERR] transaction prepare statement err : ", err)
			return err
		}
		defer stmt.Close()

		result, err := stmt.Exec(videolink, lastid)
		if err != nil {
			log.Println("[ERR] statement execution err : ", err)
			return err
		}

		rowcnt, err := result.RowsAffected()
		if err != nil {
			log.Println("[ERR] rows affected err : ", err)
			return err
		}

		if rowcnt == 0 {
			log.Println("[LOG] None affected rows")
		}
	}

	{
		if betafd != nil {
			betalink, err = projectfile.SaveBetaFile(ctx, requploadprojectinfo.ID, lastid, betafd)
			if err != nil {
				log.Println("[ERR] failed to save beta file  err : ", err)
				err = gcloud.DeleteOriginFile(ctx, requploadprojectinfo.ID, lastid, p.DB)
				if err != nil {
					log.Println("[ERR] failed to delete origin file err : ", err)
					return err
				}
				err = gcloud.DeleteVideoFile(ctx, requploadprojectinfo.ID, lastid, p.DB)
				if err != nil {
					log.Println("[ERR] failed to delete video file err : ", err)
					return err
				}
				return err
			}

			stmt, err := tx.Prepare(`UPDATE project SET beta = ?, beta_link = ? WEHRE id = ? `)
			if err != nil {
				log.Println("[ERR] transaction prepare statement err : ", err)
				return err
			}
			defer stmt.Close()

			result, err := stmt.Exec(true, betalink, lastid)
			if err != nil {
				log.Println("[ERR] statement execution err : ", err)
				return err
			}

			rowcnt, err := result.RowsAffected()
			if err != nil {
				log.Println("[ERR] rows affected err : ", err)
				return err
			}

			if rowcnt == 0 {
				log.Println("[LOG] None affected rows")
			}
		}
	}

	{

		projectimagelinks, err = image.SaveProjectImages(ctx, requploadprojectinfo.ID, lastid, imagefdarr)
		if err != nil {
			log.Println("[ERR] failed to save project images err : ", err)
			err = gcloud.DeleteOriginFile(ctx, requploadprojectinfo.ID, lastid, p.DB)
			if err != nil {
				log.Println("[ERR] failed to delete origin file err : ", err)
				return err
			}
			err = gcloud.DeleteVideoFile(ctx, requploadprojectinfo.ID, lastid, p.DB)
			if err != nil {
				log.Println("[ERR] failed to delete video file err : ", err)
				return err
			}
			err = gcloud.DeleteBetaFile(ctx, requploadprojectinfo.ID, lastid, p.DB)
			if err != nil {
				log.Println("[ERR] failed to delete beta file err : ", err)
				return err
			}
			return err
		}

		for _, projectimagelink := range projectimagelinks {
			stmt, err := tx.Prepare(`INSERT INTO image (project_id,link) VALUES(?,?) ON DUPLICATE KEY UPDATE project_id= ?, link=?`)
			if err != nil {
				log.Println("[ERR] prepare statement err : ", err)
				return err
			}

			result, err := stmt.Exec(lastid, projectimagelink, lastid, projectimagelink)
			if err != nil {
				log.Println("[ERR] statement execution err : ", err)
				return err
			}

			rowcnt, err := result.RowsAffected()
			if err != nil {
				log.Println("[ERR] rows affected err : ", err)
				return err
			}

			if rowcnt == 0 {
				log.Println("[LOG] None affected rows")
			}
			stmt.Close()
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("[ERR]  transaction commit err : ", err)
		return err
	}

	return nil

}

//GetProjectInfo is function to get project informaiton uploaded.
func GetProjectInfo(c *gin.Context) (*ReqUploadProjectInfo, *os.File, []*os.File, *os.File, *os.File, error) {

	var requploadprojectinfo *ReqUploadProjectInfo
	var imagefdarr []*os.File
	var videofd, betafd, originfd *os.File

	multipartreader, err := c.Request.MultipartReader()
	if err != nil {
		log.Println("[ERR] multipart reader err : ", err)
		return nil, nil, nil, nil, nil, err
	}
	for {
		part, err := multipartreader.NextPart()
		if err != nil {
			log.Println("[ERR] next part err : ", err)
			return nil, nil, nil, nil, nil, err
		}

		data, err := ioutil.ReadAll(part)

		switch part.Header.Get("Content-ID") {
		case "metadata":
			err = json.Unmarshal(data, requploadprojectinfo)
			if err != nil {
				log.Println("[ERR] json unmarshal err : ", err)
				return nil, nil, nil, nil, nil, err
			}
		case "image":
			file, err := os.Create(part.Header.Get("Content-Filename"))
			if err != nil {
				log.Println("[ERR] os create err : ", err)
				return nil, nil, nil, nil, nil, err
			}

			_, err = file.Write(data)
			if err != nil {
				log.Println("[ERR] file write err : ", err)
				return nil, nil, nil, nil, nil, err
			}

			imagefdarr = append(imagefdarr, file)
			file.Close()

		case "video":
			videofd, err = os.Create(part.Header.Get("Content-Filename"))
			if err != nil {
				log.Println("[ERR] os create err : ", err)
				return nil, nil, nil, nil, nil, err
			}
			_, err = videofd.Write(data)
			if err != nil {
				log.Println("[ERR] file write err : ", err)
				return nil, nil, nil, nil, nil, err
			}

			videofd.Close()

		case "beta":
			betafd, err = os.Create(part.Header.Get("Content-Filename"))
			if err != nil {
				log.Println("[ERR] os create err : ", err)
				return nil, nil, nil, nil, nil, err
			}
			_, err = betafd.Write(data)
			if err != nil {
				log.Println("[ERR] file write err : ", err)
				return nil, nil, nil, nil, nil, err
			}
			betafd.Close()

		case "origin":
			originfd, err = os.Create(part.Header.Get("Content-Filename"))
			if err != nil {
				log.Println("[ERR] os create err : ", err)
				return nil, nil, nil, nil, nil, err
			}
			_, err = originfd.Write(data)
			if err != nil {
				log.Println("[ERR] file write err : ", err)
				return nil, nil, nil, nil, nil, err
			}
			originfd.Close()
		}
	}

	return requploadprojectinfo, videofd, imagefdarr, betafd, originfd, nil
}

func (p *Pjt) InsertProjectInfo(requploadprojectinfo *ReqUploadProjectInfo) error {

	tx, err := p.DB.Begin()
	if err != nil {
		log.Println("[ERR] begin err : ", err)
		return err
	}

	defer tx.Rollback()

	tx.Prepare(`INSERT INTO project (title,category_id,description,price,video_link,`)

	return nil
}

func (p *Pjt) UpdateProjectInfo(c *gin.Context) error {
	var videoobject gcloud.Video
	var image gcloud.Image
	var projectimagelinks []string
	var videolink, betalink, originlink string
	var projectfile gcloud.ProjectFile

	requploadprojectinfo, videofd, imagefdarr, betafd, originfd, err := GetProjectInfo(c)
	if err != nil {
		log.Println("[ERR] failed to get project informaiton err : ", err)
		return err
	}

	ctx := appengine.NewContext(c.Request)

	tx, err := p.DB.Begin()
	if err != nil {
		log.Println("[ERR] begin err : ", err)
		return err
	}

	defer DeleteAllFiles(ctx, requploadprojectinfo.ID, requploadprojectinfo.ProjectID, p.DB)
	defer tx.Rollback()

	{
		stmt, err := tx.Prepare(`UPDATE project SET category_id = ?, title = ?, description = ?, price = ?,updated_at) VALUES(?,?,?,?,NOW() WHERE user_id = ? AND id = ?)`)
		if err != nil {
			log.Println("[ERR] prepare statement err : ", err)
			return err
		}
		defer stmt.Close()

		result, err := stmt.Exec(requploadprojectinfo.Category, requploadprojectinfo.Title, requploadprojectinfo.Desc, requploadprojectinfo.Price, requploadprojectinfo.ID, requploadprojectinfo.ProjectID)
		if err != nil {
			log.Println("[ERR] stmt exec err : ", err)
			return err
		}

		rowcnt, err := result.RowsAffected()
		if err != nil {
			log.Println("[ERR] rows affected err : ", err)
			return err
		}

		if rowcnt == 0 {
			log.Println("[LOG] None affected rows")
		}
	}

	{
		if originfd != nil {
			originlink, err = projectfile.SaveOriginFile(ctx, requploadprojectinfo.ID, requploadprojectinfo.ProjectID, originfd)
			if err != nil {
				log.Println("[ERR] failed to save origin file err : ", err)
				return err
			}
			stmt, err := tx.Prepare(`UPDATE project SET origin_link = ? WHERE id = ?`)
			if err != nil {
				log.Println("[ERR]prepare statement err : ", err)
				return err
			}
			defer stmt.Close()

			result, err := stmt.Exec(originlink, requploadprojectinfo.ID)
			if err != nil {
				log.Println("[ERR] stmt exec err ; ", err)
				return err
			}

			rowcnt, err := result.RowsAffected()
			if err != nil {
				log.Println("[ERR] rows affected err : ", err)
				return err
			}

			if rowcnt == 0 {
				log.Println("[LOG] None affected rows")
			}
		}
	}

	{
		if videofd != nil {
			videolink, err = videoobject.SaveVideoFile(ctx, requploadprojectinfo.ID, requploadprojectinfo.ProjectID, videofd)
			if err != nil {
				log.Println("[ERR] failed to save video file err : ", err)
				return err
			}
		}

		stmt, err := tx.Prepare(`UPDATE project SET video_link = ? WHERE id = ? `)
		if err != nil {
			log.Println("[ERR] transaction prepare statement err : ", err)
			return err
		}
		defer stmt.Close()

		result, err := stmt.Exec(videolink, requploadprojectinfo.ID)
		if err != nil {
			log.Println("[ERR] statement execution err : ", err)
			return err
		}

		rowcnt, err := result.RowsAffected()
		if err != nil {
			log.Println("[ERR] rows affected err : ", err)
			return err
		}

		if rowcnt == 0 {
			log.Println("[LOG] None affected rows")
		}
	}

	{
		if betafd != nil {
			betalink, err = projectfile.SaveBetaFile(ctx, requploadprojectinfo.ID, requploadprojectinfo.ProjectID, betafd)
			if err != nil {
				log.Println("[ERR] failed to save beta file  err : ", err)
				return err
			}

			stmt, err := tx.Prepare(`UPDATE project SET beta = ?, beta_link = ? WEHRE id = ? `)
			if err != nil {
				log.Println("[ERR] transaction prepare statement err : ", err)
				return err
			}
			defer stmt.Close()

			result, err := stmt.Exec(true, betalink, requploadprojectinfo.ID)
			if err != nil {
				log.Println("[ERR] statement execution err : ", err)
				return err
			}

			rowcnt, err := result.RowsAffected()
			if err != nil {
				log.Println("[ERR] rows affected err : ", err)
				return err
			}

			if rowcnt == 0 {
				log.Println("[LOG] None affected rows")
			}
		}
	}

	{
		if imagefdarr != nil {
			projectimagelinks, err = image.SaveProjectImages(ctx, requploadprojectinfo.ID, requploadprojectinfo.ProjectID, imagefdarr)
			if err != nil {
				log.Println("[ERR] failed to save project images err : ", err)
				return err
			}

			for _, projectimagelink := range projectimagelinks {
				stmt, err := tx.Prepare(`INSERT INTO image (project_id,link) VALUES(?,?) ON DUPLICATE KEY UPDATE project_id= ?, link=?`)
				if err != nil {
					log.Println("[ERR] prepare statement err : ", err)
					return err
				}

				result, err := stmt.Exec(requploadprojectinfo.ID, projectimagelink, requploadprojectinfo.ID, projectimagelink)
				if err != nil {
					log.Println("[ERR] statement execution err : ", err)
					return err
				}

				rowcnt, err := result.RowsAffected()
				if err != nil {
					log.Println("[ERR] rows affected err : ", err)
					return err
				}

				if rowcnt == 0 {
					log.Println("[LOG] None affected rows")
				}
				stmt.Close()
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("[ERR] transacntion commit err : ", err)
		return err
	}

	return nil
}

//DeleteAllFiles is fuction to delete files(image,beta, origin, etc) saved already
func DeleteAllFiles(ctx context.Context, userid int64, projectid int64, database *sql.DB) error {
	err := gcloud.DeleteVideoFile(ctx, userid, projectid, database)
	if err != nil {
		log.Println("[ERR] failed to delete video file err : ", err)
		return err
	}

	err = gcloud.DeleteBetaFile(ctx, userid, projectid, database)
	if err != nil {
		log.Println("[ERR] failed to delete beta file err : ", err)
		return err
	}

	err = gcloud.DeleteOriginFile(ctx, userid, projectid, database)
	if err != nil {
		log.Println("[ERR] failed to delete original file err : ", err)
		return err
	}

	err = gcloud.DeleteProjectImages(ctx, userid, projectid, database)
	if err != nil {
		log.Println("[ERR] failed to delete project images err : ", err)
		return err
	}

	return nil
}
