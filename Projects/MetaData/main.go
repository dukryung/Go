package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	database, err := sql.Open("mysql", "dukryung:superunderdog@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Println("[ERR] database open err : ", err)
		panic(err)
	}

	_, err = database.Exec("CREATE DATABASE IF NOT EXISTS " + "sideproject" + ";")
	if err != nil {
		log.Println("[LOG] database exec err : ", err)
	}

	_, err = database.Exec("USE " + "sideproject")
	if err != nil {
		log.Println("[LOG] database exec err : ", err)
	}

	err = InsertUserInfo(database)
	if err != nil {
		log.Println("[ERR] insert user info err : ", err)
		panic(err)
	}

	err = InsertProjectInfo(database)
	if err != nil {
		log.Println("[ERR] insert project info err : ", err)
		panic(err)
	}

	err = InsertProjectRankInfo(database)
	if err != nil {
		log.Println("[ERR] insert project info err : ", err)
		panic(err)
	}
}

//InsertUserInfo is function to insert user information.
func InsertUserInfo(database *sql.DB) error {
	for i := 1; i < 1001; i++ {

		sessionid := fmt.Sprintf("sessionid_%d", i)
		name := fmt.Sprintf("dukryung_%d", i)
		nickname := fmt.Sprintf("duck_%d", i)
		email := fmt.Sprintf("dukryung_%d@naver.com", i)
		agreemarketing := true
		imagelink := fmt.Sprintf("dukryung/%d", i)
		introduction := fmt.Sprintf("introduction_%d", i)
		social := true
		cash := i * 100
		bank := 123
		account := fmt.Sprintf("account_%d", i)
		agreepolicy := true

		query := fmt.Sprintf(`INSERT IGNORE INTO user 
		(session_id, name, nickname, email, agree_email_marketing,image_link,introduction,social,cash,bank,account,agree_policy, created_at, updated_at) 
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?, NOW(), NOW())`)

		stmt, err := database.Prepare(query)
		if err != nil {
			log.Println("[ERR] prepare statement err : ", err)
			panic(err)
		}
		defer stmt.Close()

		result, err := stmt.Exec(sessionid, name, nickname, email, agreemarketing, imagelink, introduction, social, cash, bank, account, agreepolicy)
		if err != nil {
			log.Println("[ERR] statement execution err : ", err)
			panic(err)
		}

		affectedrowscnt, err := result.RowsAffected()
		if err != nil {
			log.Println("[ERR] affected rows err : ", err)
			panic(err)
		}

		log.Println("[LOG] affected rows count : ", affectedrowscnt)
	}

	return nil
}

func InsertProjectInfo(database *sql.DB) error {
	for i := 1; i < 1001; i++ {

		for j := 1; j < 6; j++ {
			userid := i
			categoryid := 1
			title := fmt.Sprintf("title_%d_%d", i, j)
			description := fmt.Sprintf("description_%d_%d", i, j)
			price := j * 1000
			sellcount := j * 2
			totalupvotecount := j * 10
			commentcount := j * 10
			videolink := fmt.Sprintf("videolink/%d/%d", i, j)
			beta := true
			betalink := fmt.Sprintf("beta/%d/%d", i, j)
			originlink := fmt.Sprintf("origin/%d/%d", i, j)

			query := fmt.Sprintf(`INSERT IGNORE INTO project 
			(user_id, category_id, title, description, price,sell_count,total_upvote_count,comment_count,video_link,beta,beta_link,origin_link, created_at, updated_at) 
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?, NOW(), NOW())`)

			stmt, err := database.Prepare(query)
			if err != nil {
				log.Println("[ERR] prepare statement err : ", err)
				panic(err)
			}
			defer stmt.Close()

			result, err := stmt.Exec(userid, categoryid, title, description, price, sellcount, totalupvotecount, commentcount, videolink, beta, betalink, originlink)
			if err != nil {
				log.Println("[ERR] statement execution err : ", err)
				panic(err)
			}

			affectedrowscnt, err := result.RowsAffected()
			if err != nil {
				log.Println("[ERR] affected rows err : ", err)
				panic(err)
			}

			log.Println("[LOG] affected rows count : ", affectedrowscnt)

		}
	}

	return nil
}

func InsertProjectRankInfo(database *sql.DB) error {
	for i := 1; i < 1001; i++ {

		projectid := fmt.Sprintf("%d", i)
		score := fmt.Sprintf("%f", rand.Float64()*10+80)
		rank := 0

		query := fmt.Sprintf(`INSERT IGNORE INTO project_rank 
		(project_id, score, rank, created_at, updated_at) 
		VALUES(?,?,?,NOW(), NOW())`)

		stmt, err := database.Prepare(query)
		if err != nil {
			log.Println("[ERR] prepare statement err : ", err)
			panic(err)
		}
		defer stmt.Close()

		result, err := stmt.Exec(projectid, score, rank)
		if err != nil {
			log.Println("[ERR] statement execution err : ", err)
			panic(err)
		}

		affectedrowscnt, err := result.RowsAffected()
		if err != nil {
			log.Println("[ERR] affected rows err : ", err)
			panic(err)
		}

		log.Println("[LOG] affected rows count : ", affectedrowscnt)
	}
	return nil
}
