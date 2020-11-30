package main

import (
	"html/template"
	"os"
)

//User is user's lnformation
type User struct {
	Name  string
	Email string
	Age   int
}

//IsOld is check age
func (u User) IsOld() bool {
	return u.Age > 30
}

func main() {

	user := User{Name: "dukryung", Email: "dukryung@naver.com", Age: 32}
	user2 := User{Name: "aaa", Email: "aaa@naver.com", Age: 22}
	users := []User{user, user2}
	tmpl, err := template.New("tmpl1").ParseFiles("templates/tmpl1.tmpl", "templates/tmpl2.tmpl")
	if err != nil {
		panic(err)
	}

	tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", users)

}
