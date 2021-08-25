package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type mariadbHandler struct {
	db *sql.DB
}

func MakeDBHandler(databasename string) *sql.DB {
	return NewMariaDBHandler(databasename)
}

func NewMariaDBHandler(databasename string) *sql.DB {

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
		phonenum     VARCHAR(50),
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

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS salelist (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		user_id BIGINT UNSIGNED,
		project_id BIGINT UNSIGNED
		created_at 		   TIMESTAMP,
		UNIQUE INDEX idx_id (id),
		INDEX idx_created_at (created_at)
  		);
		`)

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
		impuid	 		TEXT,
		merchant_id 	TEXT,
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
		impuid	 		TEXT,
		merchant_id 	TEXT,
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

	stmt, err = database.Prepare(`CREATE TABLE IF NOT EXISTS payment (
		id BIGINT  UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		impuid VARCHAR(20) DEFAULT 'imp53094657',
		pg VARCHAR(20) DEFAULT 'html5_inicis',
		pay_method VARCHAR(10) DEFAULT 'card',
		created_at TIMESTAMP,
		UNIQUE INDEX idx_id (id),
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

	return database
}
