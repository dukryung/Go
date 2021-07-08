package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var datesub = "DATE_SUB(NOW(), INTERVAL 1 DAY)"

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

	err = InsertUserRankInfo(database)
	if err != nil {
		log.Println("[ERR] insert user rank info err : ", err)
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

	err = InsertBuyHistory(database)
	if err != nil {
		log.Println("[ERR] insert buy history info err : ", err)
		panic(err)
	}

	err = InsertSellHistory(database)
	if err != nil {
		log.Println("[ERR] insert buy history info err : ", err)
		panic(err)
	}

	err = InsertWithdrawHistory((database))
	if err != nil {
		log.Println("[ERR] insert withdraw history info err : ", err)
		panic(err)
	}
}

//InsertUserInfo is function to insert user information.
func InsertUserInfo(database *sql.DB) error {
	n := 1
	date := time.Now().AddDate(0, 0, -9)

	for i := 1; i < 1001; i++ {

		if i == (100*n + 1) {
			date = date.AddDate(0, 0, 1)
			n++
		}

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
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)

		stmt, err := database.Prepare(query)
		if err != nil {
			log.Println("[ERR] prepare statement err : ", err)
			panic(err)
		}
		defer stmt.Close()

		result, err := stmt.Exec(sessionid, name, nickname, email, agreemarketing, imagelink, introduction, social, cash, bank, account, agreepolicy, date, date)
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

func InsertUserRankInfo(database *sql.DB) error {
	n := 1
	date := time.Now().AddDate(0, 0, -9)

	for i := 1; i < 1001; i++ {

		if i == (100*n + 1) {
			date = date.AddDate(0, 0, 1)
			n++
		}

		userid := fmt.Sprintf("%d", i)
		score := fmt.Sprintf("%f", rand.Float64()*10+80)

		query := fmt.Sprintf(`INSERT IGNORE INTO user_rank 
		(user_id, score, created_at, updated_at) 
		VALUES(?,?,?,?)`)

		stmt, err := database.Prepare(query)
		if err != nil {
			log.Println("[ERR] prepare statement err : ", err)
			panic(err)
		}
		defer stmt.Close()

		result, err := stmt.Exec(userid, score, date, date)
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
	n := 1
	date := time.Now().AddDate(0, 0, -9)
	for i := 1; i < 1001; i++ {
		if i == (100*n + 1) {
			date = date.AddDate(0, 0, 1)
			n++
		}
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
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)

			stmt, err := database.Prepare(query)
			if err != nil {
				log.Println("[ERR] prepare statement err : ", err)
				panic(err)
			}
			defer stmt.Close()

			result, err := stmt.Exec(userid, categoryid, title, description, price, sellcount, totalupvotecount, commentcount, videolink, beta, betalink, originlink, date, date)
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
	n := 1
	var projectid int
	date := time.Now().AddDate(0, 0, -9)
	for i := 1; i < 1001; i++ {
		if i == (100*n + 1) {
			date = date.AddDate(0, 0, 1)
			n++
		}
		for j := 1; j < 6; j++ {
			projectid++
			score := fmt.Sprintf("%f", rand.Float64()*10+80)

			query := fmt.Sprintf(`INSERT IGNORE INTO project_rank 
		(project_id, score, created_at, updated_at) 
		VALUES(?,?,?,?)`)

			stmt, err := database.Prepare(query)
			if err != nil {
				log.Println("[ERR] prepare statement err : ", err)
				panic(err)
			}
			defer stmt.Close()

			result, err := stmt.Exec(projectid, score, date, date)
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

func InsertBuyHistory(database *sql.DB) error {
	var n = 1
	var sellerid int
	var userid int
	var sellernickname string
	for i := 1; i < 101; i++ {

		for j := 1; j < 6; j++ {
			projectid := j
			price := 10000 + i
			importid := 1000 + i

			if i <= 5*n {
				userid = 80 + i
				sellerid = n
				sellernickname = fmt.Sprintf("dukryung_%d", n)
			} else {
				n++
				sellerid = n
				sellernickname = fmt.Sprintf("dukryung_%d", n)
			}

			stmt, err := database.Prepare(`INSERT IGNORE INTO buy_history
							  (user_id,project_id,seller_id,seller_nickname,price,import_id)
							  VALUES(?,?,?,?,?,?)`)
			if err != nil {
				log.Println("[ERR] prepare statement err : ", err)
				return err
			}

			defer stmt.Close()

			result, err := stmt.Exec(userid, projectid, sellerid, sellernickname, price, importid)
			if err != nil {
				log.Println("[ERR] stmt excution err : ", err)
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
func InsertSellHistory(database *sql.DB) error {
	var n = 1
	var sellerid int
	var userid int
	var sellernickname string
	for i := 1; i < 101; i++ {

		for j := 1; j < 6; j++ {
			projectid := i
			price := 10000 + i
			importid := 1000 + i

			if i <= 5*n {
				userid = 80 + i
				sellerid = n
				sellernickname = fmt.Sprintf("dukryung_%d", n)
			} else {
				n++
				sellerid = n
				sellernickname = fmt.Sprintf("dukryung_%d", n)
			}

			stmt, err := database.Prepare(`INSERT IGNORE INTO sell_history
							  (user_id,project_id,buyer_id,buyer_nickname,price,import_id)
							  VALUES(?,?,?,?,?,?)`)
			if err != nil {
				log.Println("[ERR] prepare statement err : ", err)
				return err
			}

			defer stmt.Close()

			result, err := stmt.Exec(userid, projectid, sellerid, sellernickname, price, importid)
			if err != nil {
				log.Println("[ERR] stmt excution err : ", err)
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

func InsertWithdrawHistory(database *sql.DB) error {
	var n = 1

	var userid int

	for i := 1; i < 201; i++ {

		cash := 10000 + i
		if i < 8*n {
			userid = n + 90
		} else {
			n++
			userid = n + 90
		}

		stmt, err := database.Prepare(`INSERT IGNORE INTO withdraw_history
							  (user_id,amount)
							  VALUES(?,?)`)
		if err != nil {
			log.Println("[ERR] prepare statement err : ", err)
			return err
		}

		defer stmt.Close()

		result, err := stmt.Exec(userid, cash)
		if err != nil {
			log.Println("[ERR] stmt excution err : ", err)
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
