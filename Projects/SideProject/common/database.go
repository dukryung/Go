package common

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type DBHandler interface {
	ReadProjectList(*ReqProjectsOfTheDay) *ResProjectsOfTheDay
	ReadArtistList() *ResArtistOfTheMonth
	ReadUserInfo(string) *ResUserInfo
	CreateUserInfo(AuthUserInfo) error
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
		email 	     VARCHAR(200),
		image_link   TEXT,
		introduction VARCHAR(200),
		created_at   TIMESTAMP,
		updated_at   TIMESTAMP,
		social       TEXT,
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
  		category_id    BIGINT UNSIGNED,
  		title 			   VARCHAR(200),
  		description 	   TEXT,
  		price 			   INT,
  		sell_count 		   INT,
  		total_upvote_count INT,
  		comment_count 	   INT,
  		beta 			   BOOLEAN,
  		created_at 		   TIMESTAMP,
  		updated_at 		   TIMESTAMP,
  		UNIQUE INDEX idx_id (id),
		INDEX idx_user_id (user_id),
		INDEX idx_created_id (created_id)
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

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS video (
  		id 			BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		project_id  BIGINT UNSIGNED,
  		link 		TEXT,
  		created_at  TIMESTAMP,
  		updated_at  TIMESTAMP,
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
  		user_id BIGINT UNSIGNED,
  		project_id BIGINT UNSIGNED,
  		seller_id BIGINT UNSIGNED,
  		import_id int,
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

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS sell_history (
  		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		user_id BIGINT UNSIGNED,
  		project_id BIGINT UNSIGNED,
  		buyer_id BIGINT UNSIGNED,
  		import_id int,
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

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS withdraw_history (
  		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		user_id BIGINT UNSIGNED,
  		amount  INT,
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

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS account (
  		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		user_id BIGINT UNSIGNED,
  		project_id BIGINT UNSIGNED,
  		cash  INT,
  		bank  INT,
  		account TEXT,
  		created_at TIMESTAMP,
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

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS withdraw_status (
  		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		user_id BIGINT UNSIGNED,
  		cash  INT,
  		created_at TIMESTAMP,
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
  		user_id BIGINT UNSIGNED,
  		project_id BIGINT UNSIGNED,
  		text TEXT,
  		created_at TIMESTAMP,
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

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS reply (
  		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  		user_id BIGINT UNSIGNED,
  		comment_id BIGINT UNSIGNED,
  		text TEXT,
  		created_at TIMESTAMP,
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

	return &mariadbHandler{db: database}
}

func (m *mariadbHandler) ReadProjectList(reqpod *ReqProjectsOfTheDay) *ResProjectsOfTheDay {

	var respod *ResProjectsOfTheDay
	respod = &ResProjectsOfTheDay{Date: reqpod.DemandDate}

	period, err := strconv.Atoi(reqpod.DemandPeriod)
	if err != nil {
		log.Println("[LOG] rows err : ", err)
		return nil
	}

	stmt, err := m.db.Prepare(`SELECT 
						p.id, 
						p.title, 
						c.code,
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
									INNER JOIN category AS c ON p.category_id = c.id
				  WHERE p_r.created_at BETWEEN DATE_FORMAT(DATE_SUB( ? ,INTERVAL ? DAY),"%Y-%m-%d") AND DATE_FORMAT(DATE_ADD( ? ,INTERVAL ? DAY), "%Y-%m-%d")
				  LIMIT 10;`)

	if err != nil {
		log.Println("[LOG] stmt err : ", err)
		return nil
	}
	defer stmt.Close()

	var rows *sql.Rows

	if period == 1 {
		rows, err = stmt.Query(reqpod.DemandDate, period-1, reqpod.DemandDate, period)
		if err != nil {
			log.Println("[LOG] rows err : ", err)
			return nil
		}
	} else if period == 7 || period == 30 {
		rows, err = stmt.Query(reqpod.DemandDate, period-1, reqpod.DemandDate, 1)
		if err != nil {
			log.Println("[LOG] rows err : ", err)
			return nil
		}
	} else {
		return nil
	}

	defer rows.Close()

	var project ProjectList
	respod.Date = reqpod.DemandDate
	respod.Period = reqpod.DemandPeriod
	respod.Total = "0"
	respod.RankLastNumber = "0"

	//var projectid, ranking uint64
	var projectid, ranking, title, categorycode, desc, createdat, sellcount, nickname, commentcount, totalupvotecount, price, beta string
	for rows.Next() {
		err := rows.Scan(&projectid, &title, &categorycode, &desc, &createdat, &sellcount, &nickname, &commentcount, &totalupvotecount, &price, &beta, &ranking)
		if err != nil {
			log.Println("[LOG] scan err : ", err)
			return nil
		}
		//project.ID = strconv.FormatUint(projectid, 10)
		//project.Rank = strconv.FormatUint(ranking, 10)
		project.ID = projectid
		project.Title = title
		project.CategoryCode = categorycode
		project.Description = desc
		project.CreatedAt = createdat
		project.SellCount = sellcount
		project.UserNickName = nickname
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
			log.Println("[LOG] stmt query err : ", err)
			return nil
		}

		for rows.Next() {
			err := rows.Scan(&link)
			if err != nil {
				log.Println("[LOG] rows scan err : ", err)
				return nil
			}
			respod.Project[i].ImageLink = link
		}
	}

	stmt, err = m.db.Prepare(`SELECT COUNT(*) 
				  FROM project 
				  WHERE created_at 
				  BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL ? DAY),"%Y-%m-%d") AND DATE_FORMAT(DATE_ADD(?, INTERVAL ? DAY),"%Y-%m-%d")`)

	if period == 1 {
		rows, err = stmt.Query(reqpod.DemandDate, period-1, reqpod.DemandDate, period)
		if err != nil {
			log.Println("[LOG] rows err : ", err)
			return nil
		}
	} else if period == 7 || period == 30 {
		rows, err = stmt.Query(reqpod.DemandDate, period-1, reqpod.DemandDate, 1)
		if err != nil {
			log.Println("[LOG] rows err : ", err)
			return nil
		}
	} else {
		return nil
	}

	var projectcnt string
	for rows.Next() {
		err := rows.Scan(&projectcnt)
		if err != nil {
			return nil
		}
		respod.Total = projectcnt
	}

	return respod
}

func (m *mariadbHandler) ReadArtistList() *ResArtistOfTheMonth {

	var resaom *ResArtistOfTheMonth
	resaom = &ResArtistOfTheMonth{}

	stmt, err := m.db.Prepare(`SELECT 
								u.nickname,
								u.introduction,
								u.image_link,
								RANK() OVER(ORDER BY u_r.score DESC)
								FROM user AS u INNER JOIN user_rank AS u_r ON u.id = u_r.user_id 
								WHERE u_r.created_at BETWEEN  DATE_FORMAT(DATE_SUB(NOW(),INTERVAL 1 MONTH),"%Y-%m-%d") AND DATE_FORMAT(NOW(),"%Y-%m-%d")
								LIMIT 5;`)

	if err != nil {
		log.Println("[LOG]prepare statement err : ", err)
		return nil
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Println("[LOG] query err : ", err)
		return nil
	}

	defer rows.Close()

	var artist ArtistList

	for rows.Next() {
		//rows.Scan(&artist.NickName, &artist.Introduction, &artist.ImageLink)
		rows.Scan(&artist.NickName, &artist.Introduction, &artist.ImageLink, &artist.Rank)
		resaom.Artist = append(resaom.Artist, artist)
	}

	return resaom
}

func (m *mariadbHandler) ReadUserInfo(sessionid string) *ResUserInfo {
	var resuser *ResUserInfo

	resuser = &ResUserInfo{}

	stmt, err := m.db.Prepare(`SELECT id, nickname
				  FROM user WHERE session_id = ?`)
	if err != nil {
		log.Println("[ERR] prepare stmt err : ", err)
		return nil
	}

	defer stmt.Close()

	rows, err := stmt.Query(sessionid)
	if err != nil {
		log.Println("[ERR] stmt query err : ", err)
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&resuser.ID, &resuser.NickName)
		if err != nil {
			log.Println("[ERR] stmt query err : ", err)
			return nil
		}
	}

	return resuser
}

func (m *mariadbHandler) CreateUserInfo(authuserinfo AuthUserInfo) error {
	tx, err := m.db.Begin()
	if err != nil {
		log.Println("[ERR] transaction begin err : ", err)
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT IGNORE INTO user (session_id,email,social) VALUES(?,?,?)`)
	if err != nil {
		log.Println("[ERR]  prepared statement err : ", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(authuserinfo.ID, authuserinfo.Email, authuserinfo.Social)
	if err != nil {
		log.Println("[ERR] exec statement err : ", err)
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

	err = tx.Commit()
	if err != nil {
		log.Println("[ERR] transacntion commit err : ", err)
		return err
	}

	return nil
}

func (m *mariadbHandler) UpdateUserInfo(sessionid string, reqjoininfo ReqJoinInfo) error {

	tx, err := m.db.Begin()
	if err != nil {
		log.Println("[ERR] begin err : ", err)
		return err
	}

	stmt, err := tx.Prepare(`UPDATE user SET name=?, nickname=?,image_link=?,introduction=?,social=?,updated_at=NOW()
	WHERE session_id=?`)

	defer stmt.Close()

	stmt.Exec(reqjoininfo.UserInfo.Name, reqjoininfo.UserInfo.Nickname)

	return nil
}
