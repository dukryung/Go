package toml

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Response struct {
	Code         string             `json:"code"`
	Message      string             `json:"message"`
	Detail       string             `json:"detail"`
	Path         string             `json:"path"`
	Movie        []Movieinfo        `json:"movie"`
	Jobstatement []Jobstatementinfo `json:"jobstatement"`
}

type Movieinfo struct {
	Href     string `json:"href"`
	Name     string `json:"name"`
	Movie_id string `json:"movie_id"`
}
type Jobstatementinfo struct {
	Href     string `json:"href"`
	Movie_id string `json:"movie_id"`
}

func Gettoml() {
	var Res Response
	_, err := toml.DecodeFile("./test.toml", &Res)
	if err != nil {
		fmt.Println("[ERR] toml DecodeFile", err)
	}
	fmt.Println("[LOG] Res :", Res)

}
