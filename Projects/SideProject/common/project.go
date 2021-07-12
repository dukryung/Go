package common

import (
	"context"
	"database/sql"
	"os"
	"time"

	"google.golang.org/appengine"

	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

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
	Title       string `json:"title"`
	Category    string `json:"category"`
	Desc        string `json:"desc"`
	Price       int    `json:"price"`
	YoutubeLink string `json:"youtube"`
}

type project struct {
	db DBHandler
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

func (m *mariadbHandler) ReadProjectDetailArtistProjectInfo(reqprojectdetailinfo *ReqProjectDetailInfo) (*ResProjectDetailInfo, error) {

	var resprojectdetailinfo *ResProjectDetailInfo

	// -------{TODO : delete rank column and test query  in database's table.}
	stmt, err := m.db.Prepare(`SELECT 
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

func (m *mariadbHandler) ReadProjectDetailArtistProjectImagesInfo(c *gin.Context) (*ResProjectDetailInfo, error) {
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

	stmt, err := m.db.Prepare(`SELECT 
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

func (m *mariadbHandler) ReadProjectDetailCommentInfo(c *gin.Context) (*ResProjectDetailInfo, error) {
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

	stmt, err := m.db.Prepare(`SELECT
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

func (m *mariadbHandler) CreateProjectInfo(c *gin.Context, userid *int) error {
	var isbeta = false
	var videoobject Video
	var projectfile ProjectFile
	var image Image
	var lastid int64
	var projectimagelinks []string
	var videolink, betalink, originlink string

	requploadprojectinfo, videofd, imagefdarr, betafd, originfd, err := GetProjectInfo(c)
	if err != nil {
		log.Println("[ERR] failed to get project informaiton err : ", err)
		return err
	}

	ctx := appengine.NewContext(c.Request)

	tx, err := m.db.Begin()
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

		result, err := stmt.Exec(userid, requploadprojectinfo.Category, requploadprojectinfo.Title, requploadprojectinfo.Desc, requploadprojectinfo.Price, requploadprojectinfo.YoutubeLink, isbeta, betalink)
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
		originlink, err = projectfile.SaveOriginFile(ctx, userid, lastid, originfd)
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
			videolink, err = videoobject.SaveVideoFile(ctx, userid, lastid, videofd)
			if err != nil {
				log.Println("[ERR] failed to save video file  err : ", err)
				err = DeleteOriginFile(ctx, userid, lastid, m.db)
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
			betalink, err = projectfile.SaveBetaFile(ctx, userid, lastid, betafd)
			if err != nil {
				log.Println("[ERR] failed to save beta file  err : ", err)
				err = DeleteOriginFile(ctx, userid, lastid, m.db)
				if err != nil {
					log.Println("[ERR] failed to delete origin file err : ", err)
					return err
				}
				err = DeleteVideoFile(ctx, userid, lastid, m.db)
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

		projectimagelinks, err = image.SaveProjectImages(ctx, userid, lastid, imagefdarr)
		if err != nil {
			log.Println("[ERR] failed to save project images err : ", err)
			err = DeleteOriginFile(ctx, userid, lastid, m.db)
			if err != nil {
				log.Println("[ERR] failed to delete origin file err : ", err)
				return err
			}
			err = DeleteVideoFile(ctx, userid, lastid, m.db)
			if err != nil {
				log.Println("[ERR] failed to delete video file err : ", err)
				return err
			}
			err = DeleteBetaFile(ctx, userid, lastid, m.db)
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

func (m *mariadbHandler) InsertProjectInfo(requploadprojectinfo *ReqUploadProjectInfo) error {

	tx, err := m.db.Begin()
	if err != nil {
		log.Println("[ERR] begin err : ", err)
		return err
	}

	defer tx.Rollback()

	tx.Prepare(`INSERT INTO project (title,category_id,description,price,video_link,`)

	return nil
}

func (m *mariadbHandler) UpdateProjectInfo(c *gin.Context, userid *int) error {
	var videoobject Video
	var image Image
	var projectimagelinks []string
	var videolink, betalink, originlink string
	var projectfile ProjectFile

	requploadprojectinfo, videofd, imagefdarr, betafd, originfd, err := GetProjectInfo(c)
	if err != nil {
		log.Println("[ERR] failed to get project informaiton err : ", err)
		return err
	}

	ctx := appengine.NewContext(c.Request)

	tx, err := m.db.Begin()
	if err != nil {
		log.Println("[ERR] begin err : ", err)
		return err
	}

	defer DeleteAllFiles(ctx, userid, requploadprojectinfo.ID, m.db)
	defer tx.Rollback()

	{
		stmt, err := tx.Prepare(`UPDATE project SET category_id = ?, title = ?, description = ?, price = ?,updated_at) VALUES(?,?,?,?,NOW() WHERE id = ?)`)
		if err != nil {
			log.Println("[ERR] prepare statement err : ", err)
			return err
		}
		defer stmt.Close()

		result, err := stmt.Exec(requploadprojectinfo.Category, requploadprojectinfo.Title, requploadprojectinfo.Desc, requploadprojectinfo.Price, requploadprojectinfo.ID)
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
			originlink, err = projectfile.SaveOriginFile(ctx, userid, requploadprojectinfo.ID, originfd)
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
			videolink, err = videoobject.SaveVideoFile(ctx, userid, requploadprojectinfo.ID, videofd)
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
			betalink, err = projectfile.SaveBetaFile(ctx, userid, requploadprojectinfo.ID, betafd)
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
			projectimagelinks, err = image.SaveProjectImages(ctx, userid, requploadprojectinfo.ID, imagefdarr)
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
func DeleteAllFiles(ctx context.Context, userid *int, projectid int64, database *sql.DB) error {
	err := DeleteVideoFile(ctx, userid, projectid, database)
	if err != nil {
		log.Println("[ERR] failed to delete video file err : ", err)
		return err
	}

	err = DeleteBetaFile(ctx, userid, projectid, database)
	if err != nil {
		log.Println("[ERR] failed to delete beta file err : ", err)
		return err
	}

	err = DeleteOriginFile(ctx, userid, projectid, database)
	if err != nil {
		log.Println("[ERR] failed to delete original file err : ", err)
		return err
	}

	err = DeleteProjectImages(ctx, userid, projectid, database)
	if err != nil {
		log.Println("[ERR] failed to delete project images err : ", err)
		return err
	}

	return nil
}
