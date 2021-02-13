package common

import "fmt"

type userinfo struct {
	Sessionid int    `json:"sessionid`
	ID        string `json:"id"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

func (s *sqliteHandler) createUsers(userinfo *userinfo) {
	stmt, err := s.db.Prepare("INSERT INTO users (name, sessionid, createdat)  VALUES(?,?,datetime('now'))")
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		panic(err)
	}

	rowcnt, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	if rowcnt == 0 {
		fmt.Println("[LOG] anthing is not affected")
	}

}

func (s *sqliteHandler) readUsers(id int, sessionid string) []*Users {
	userlist := []*Users{}
	rows, err := s.db.Query("SELECT id, name, sessionid FROM users WHERE sessioni=?", sessionid)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var user Users
		rows.Scan(&user.ID, &user.Name, &user.SessionID)
		userlist = append(userlist, &user)
	}

	return userlist
}

func (s *sqliteHandler) updateUsers(id int, name string) {
	stmt, err := s.db.Prepare("UPDATE users SET id=?, name =?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id, name)
	if err != nil {
		panic(err)
	}

	rowcnt, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	if rowcnt == 0 {
		fmt.Println("[LOG] anthing is not affected")
	}
}

func (s *sqliteHandler) deleteUsers(name string, id int) {
	stmt, err := s.db.Prepare("DELETE FROM users where id=? AND name=?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id, name)
	if err != nil {
		panic(err)
	}

	rowcnt, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	if rowcnt == 0 {
		fmt.Println("[LOG] anything is not affected")
	}

}
