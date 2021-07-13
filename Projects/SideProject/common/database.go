package common

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type DBHandler interface {
	ReadProjectList(*ReqProjectsOfTheDay) (*ResProjectsOfTheDay, error)
	ReadArtistList() (*ResArtistOfTheMonth, error)
	ReadUserInfo(int64) *ResUserInfo
	ReadUserID(string) (*int, error)
	CreateUserInfo(AuthUserInfo) error

	SaveJoinUserInfo(*gin.Context) (string, error)
	UpdateModificationUserInfo(*ReqModificationUserInfo) error

	ReadProfileFrameInfo(int64) (*ResProfileFrameInfo, error)
	ReadProfileProjectInfo(int64) (*ResProfileProjectInfo, error)
	ReadProfileSellInfo(int64) (*ResProfileSellInfo, error)
	ReadProfileBuyInfo(int64) (*ResProfileBuyInfo, error)
	ReadProfileWithdrawInfo(int64) (*ResProfileWithdrawInfo, error)
	ReadModificationUserInfo(int64) (*ResModificationUserInfo, error)
	ReadProfileArtistInfo(int64) (*ResArtistProfileInfo, error)

	ReadPersonalInformation(int64) (*ResPersonalInformation, error)
	UpdatePersonalInformation(*ReqPersonalInformation) error

	ReadProjectDetailArtistProjectInfo(*ReqProjectDetailInfo) (*ResProjectDetailInfo, error)
	ReadProjectDetailArtistProjectImagesInfo(*gin.Context) (*ResProjectDetailInfo, error)
	ReadProjectDetailCommentInfo(*gin.Context) (*ResProjectDetailInfo, error)

	CreateProjectInfo(*gin.Context, *int) error
}

type mariadbHandler struct {
	db *sql.DB
}

//ProjectDB is structure for objectification
type ProjectDB struct {
}

//DB is interface to handle mariadb object
type DB interface {
	NewMariaDBHandler(string) *mariadbHandler
}

func MakeDBHandler(databasename string, database DB) DBHandler {
	return database.NewMariaDBHandler(databasename)
}

func (p *ProjectDB) NewMariaDBHandler(databasename string) *mariadbHandler {

	database, err := sql.Open("mysql", "dukryung:superunderdog@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Println("[LOG] database open err : ", err)
		return nil
	}

	_, err = database.Exec("CREATE DATABASE IF NOT EXISTS " + databasename + ";")
	if err != nil {
		log.Println("[LOG] database exec err : ", err)
	}

	_, err = database.Exec("USE " + databasename)
	if err != nil {
		log.Println("[LOG] database exec err : ", err)
	}

	stmt, err := database.Prepare(`CREATE TABLE IF NOT EXISTS user (
		id 		     BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		session_id   TEXT,
		name	     VARCHAR(50),
		nickname     VARCHAR(50),
		email 	     VARCHAR(320),
		agree_email_marketing BOOLEAN,
		image_link   TEXT,
		introduction VARCHAR(200),
		social       TEXT,
		cash  		 INT,
  		bank  		 INT,
  		account	     TEXT,
		agree_policy BOOLEAN,
		created_at   TIMESTAMP,
		updated_at   TIMESTAMP,
		UNIQUE INDEX idx_id (id),
		UNIQUE INDEX idx_session_id (session_id)
		);
		`)
	if err != nil {
		log.Println("[LOG] database prepare err:", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS project (
		id BIGINT 		   UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		user_id 	 	   BIGINT UNSIGNED,
  		category_id    	   BIGINT UNSIGNED,
  		title 			   VARCHAR(200),
  		description 	   TEXT,
  		price 			   INT,
  		sell_count 		   INT,
  		total_upvote_count INT,
  		comment_count 	   INT,
		video_link		   TEXT,
  		beta 			   BOOLEAN,
		beta_link		   TEXT,
		origin_link 	   TEXT,
  		created_at 		   TIMESTAMP,
  		updated_at 		   TIMESTAMP,
  		UNIQUE INDEX idx_id (id),
		INDEX idx_created_at (created_at)
  		);
  		`)

	if err != nil {
		log.Println("[LOG] database prepare err : ", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS image (
  		id BIGINT 		UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		project_id 		BIGINT UNSIGNED,
  		link 			TEXT,
  		created_at  	TIMESTAMP,
  		updated_at  	TIMESTAMP,
  		UNIQUE INDEX idx_id (id),
		INDEX idx_project_id (project_id)
  		);
  		`)

	if err != nil {
		log.Println("[LOG] database prepare err:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS beta_download_user (
  		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		user_id BIGINT UNSIGNED,
  		project_id BIGINT UNSIGNED,
  		created_at TIMESTAMP,
  		UNIQUE INDEX idx_id (id),
  		UNIQUE INDEX idx_project_id (project_id)
  		);
  		`)

	if err != nil {
		log.Println("[LOG] database prepare err:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS buy_history (
  		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		user_id 		BIGINT UNSIGNED,
  		project_id  	BIGINT UNSIGNED,
  		seller_id 		BIGINT UNSIGNED,
		seller_nickname VARCHAR(50),
		price			INT,
  		import_id int,
  		created_at TIMESTAMP,
  		UNIQUE INDEX idx_id (id)
  		);
  		`)

	if err != nil {
		log.Println("[LOG] database prepare err:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS sell_history (
  		id 				BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		user_id 		BIGINT UNSIGNED,
  		project_id 		BIGINT UNSIGNED,
  		buyer_id 		BIGINT UNSIGNED,
		buyer_nickname  VARCHAR(50),
		price			INT,
  		import_id 		INT,
  		created_at 		TIMESTAMP,
  		UNIQUE INDEX idx_id (id),
  		INDEX idx_user_id (user_id)
  		);
  		`)

	if err != nil {
		log.Println("[LOG] database prepare err:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS withdraw_history (
  		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		user_id BIGINT UNSIGNED,
  		amount  INT,
		requested_at TIMESTAMP,
		completed_at TIMESTAMP,
  		created_at TIMESTAMP,
  		UNIQUE INDEX idx_id (id),
  		INDEX idx_user_id (user_id)
  		);
  		`)

	if err != nil {
		log.Println("[LOG] database prepare err:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS withdraw_status (
  		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		user_id BIGINT UNSIGNED,
  		cash  INT,
		is_remit_req BOOLEAN,
  		created_at TIMESTAMP,
		request_at TIMESTAMP, 
  		updated_at TIMESTAMP,
  		UNIQUE INDEX idx_id (id),
  		INDEX idx_user_id (user_id)
  		);
  		`)

	if err != nil {
		log.Println("[LOG] database prepare err:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS comment (
  		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		artist_id BIGINT  UNSIGNED,
  		project_id BIGINT UNSIGNED,
  		artist_nickname   VARCHAR(50),
  		text TEXT,
  		created_at TIMESTAMP,
  		updated_at TIMESTAMP,
  		UNIQUE INDEX idx_id (id),
  		INDEX idx_artist_id (artist_id)
  		);
  		`)

	if err != nil {
		log.Println("[LOG] database prepare err:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS reply (
  		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		artist_id BIGINT UNSIGNED,
  		comment_id BIGINT UNSIGNED,
		artist_nickname VARCHAR(50),
  		text TEXT,
  		created_at TIMESTAMP,
  		updated_at TIMESTAMP,
  		UNIQUE INDEX idx_id (id),
  		INDEX idx_comment_id (comment_id)
  		);
  		`)

	if err != nil {
		log.Println("[LOG] database prepare err:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS category (
  		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		code VARCHAR(200),
  		created_at TIMESTAMP,
  		UNIQUE INDEX idx_id (id)
  		);
  		`)

	if err != nil {
		log.Println("[LOG] database prepare err:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS project_rank (
  		id BIGINT  UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		project_id BIGINT UNSIGNED,
  		score  	   FLOAT,
  		created_at TIMESTAMP,
  		updated_at TIMESTAMP,
  		UNIQUE INDEX idx_id (id),
  		INDEX idx_project_id (project_id),
		INDEX idx_created_at (created_at)
  		);
  		`)

	if err != nil {
		log.Println("[LOG] database prepare err:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS user_rank (
  		id BIGINT  UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		user_id BIGINT UNSIGNED,
  		score  	   FLOAT, 
  		created_at TIMESTAMP,
  		updated_at TIMESTAMP,
  		UNIQUE INDEX idx_id (id),
  		INDEX idx_user_id (user_id),
		INDEX idx_created_at (created_at)
  		);
  		`)

	if err != nil {
		log.Println("[LOG] database prepare err:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS upvote_status (
		id BIGINT  UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		user_id BIGINT UNSIGNED,
		project_id BIGINT UNSIGNED,
		created_at TIMESTAMP,
		UNIQUE INDEX idx_id (id),
		INDEX idx_user_id (user_id),
	 	INDEX idx_created_at (created_at)
		);
		`)

	if err != nil {
		log.Println("[LOG] database prepare err:", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("[LOG] stmt exec error : ", err)
	}

	return &mariadbHandler{db: database}
}

func (m *mariadbHandler) ReadProjectList(reqpod *ReqProjectsOfTheDay) (*ResProjectsOfTheDay, error) {

	var respod *ResProjectsOfTheDay
	respod = &ResProjectsOfTheDay{Date: reqpod.DemandDate}

	stmt, err := m.db.Prepare(`SELECT 
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
		stmt, err = m.db.Prepare(`SELECT link FROM image WHERE id = ? LIMIT 1;`)
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

	stmt, err = m.db.Prepare(`SELECT COUNT(*) 
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

func (m *mariadbHandler) ReadArtistList() (*ResArtistOfTheMonth, error) {

	var resaom *ResArtistOfTheMonth
	resaom = &ResArtistOfTheMonth{}

	stmt, err := m.db.Prepare(`SELECT 
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

	var artist Artist

	for rows.Next() {
		rows.Scan(&artist.ArtistID, &artist.NickName, &artist.Introduction, &artist.ImageLink, &artist.Rank)
		resaom.Artist = append(resaom.Artist, artist)
	}

	return resaom, nil
}
