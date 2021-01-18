package toml

import (
	"fmt"
	"time"

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

type Config struct {
	Age        int
	Cats       []string
	Pi         float64
	Perfection []int
	DOB        time.Time // requires `import time`
}

func Gettoml(filename string) {

	var tomltext = `Age = 25
	Cats = [ "Cauchy", "Plato" ]
	Pi = 3.14
	Perfection = [ 6, 28, 496, 8128 ]
	DOB = 1987-07-05T05:45:00Z`

	

	var Res Response
	_, err := toml.DecodeFile(filename, &Res)
	if err != nil {
		fmt.Println("[ERR] toml DecodeFile", err)
	}
	fmt.Println("[LOG] Res :", Res)

}
