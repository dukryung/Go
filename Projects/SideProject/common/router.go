package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/unrolled/render"

	"github.com/gorilla/mux"
)

var rd *render.Render = render.New()

const htmlIndex = `<html><body>
Logged in with <a href="/auth/google/signup">google</a>
</br>
Logged in with <a href="/auth/facebook/signup">facebook</a>
</br>
Logged in with <a href="/auth/kakao/signup">kakao</a>
</br>
Logged in with <a href="/auth/naver/signup">naver</a>
</br>
</body></html>
`

type project struct {
	db DBHandler
}

//ReqProjectsOfTheDay is structure to contain request information.
type ReqProjectsOfTheDay struct {
	DemandDate   string `json:"demand_date"`
	DemandPeriod string `json:"demand_period"`
}

//ResProjectsOfTheDay is structure to contain response information.
type ResProjectsOfTheDay struct {
	Date           string        `json:"date"`
	Project        []ProjectList `json:"project"`
	Total          string        `json:"total"`
	Period         string        `json:"period"`
	RankLastNumber string        `json:"rank_last_number"`
}

//ProjectList is structure to get project list information.
type ProjectList struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	CategoryCode string `json:"category_code"`
	Description  string `json:"desc"`
	ImageLink    string `json:"image_link"`
	CreatedAt    string `json:"created_at"`
	SellCount    string `json:"sell_cnt"`
	UserNickName string `json:"user_nickname"`
	CommentCount string `json:"comment_count"`
	UpvoteCount  string `json:"upvote_count"`
	Price        string `json:"price"`
	Beta         string `json:"beta"`
	Rank         string `json:"rank"`
}

//ResArtistOfTheMonth is sturcture to contain response artist information.
type ResArtistOfTheMonth struct {
	Artist []ArtistList `json:"artist"`
}

//ArtistList is structure to get artist list information.
type ArtistList struct {
	NickName     string `json:"nickname"`
	Introduction string `json:"introduction"`
	ImageLink    string `json:"image_link"`
	Rank         string `json:"rank"`
}

//ResUserInfo is structure to contain user information in index page.
type ResUserInfo struct {
	ID       string `json:"id"`
	NickName string `json:"nickname"`
}

type ResEmailInfo struct {
	Email string `json:"email"`
}

type ReqJoinInfo struct {
	UserInfo    User    `json:"user"`
	AccountInfo Account `json:"account"`
}
type User struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Nickname            string `json:"nickname"`
	Email               string `json:"agree_email_marketing"`
	Introduction        string `json:"introduction"`
	AgreeEmailMarketing bool   `json:"agree_email_marketing"`
}

type Account struct {
	UserID      string `json:"user_id"`
	Bank        string `json:"bank"`
	Account     string `json:"account"`
	AgreePolicy bool   `json:"agree_policy"`
}

func (p *project) GetProjectInfoHandler(w http.ResponseWriter, r *http.Request) {
	var reqpod ReqProjectsOfTheDay

	err := r.ParseForm()
	if err != nil {
		log.Println("[LOG] parseform err : ", err)
		rd.JSON(w, http.StatusInternalServerError, nil)
	}

	dmddate := r.Form.Get("demand_date")
	dmdperiod := r.Form.Get("demand_period")

	reqpod.DemandDate = dmddate
	reqpod.DemandPeriod = dmdperiod

	log.Println("[LOG] reqpod information : ", reqpod)

	respod := p.db.ReadProjectList(&reqpod)
	if respod == nil {
		rd.JSON(w, http.StatusInternalServerError, nil)
	}

	rd.JSON(w, http.StatusOK, respod)

}

func (p *project) GetArtistInfoHandler(w http.ResponseWriter, r *http.Request) {
	resaom := p.db.ReadArtistList()
	log.Println("[LOG] resaom : ", resaom)

	rd.JSON(w, http.StatusOK, resaom)
}

func (p *project) GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {

	sessionid := getSessionID(r)

	resuser := p.db.ReadUserInfo(sessionid)

	rd.JSON(w, http.StatusOK, resuser)
}

func (p *project) GetIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}

func (p *project) PutUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	sessionid := getSessionID(r)

	if sessionid == "" {
		rd.JSON(w, http.StatusUnauthorized, nil)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rd.JSON(w, http.StatusInternalServerError, nil)
	}

	var reqjoininfo ReqJoinInfo
	err = json.Unmarshal(data, &reqjoininfo)
	if err != nil {
		rd.JSON(w, http.StatusInternalServerError, nil)
	}

}

func MakeHandler(databasename string) http.Handler {
	mux := mux.NewRouter()

	pdb := &ProjectDB{}

	p := &project{db: MakeDBHandler(databasename, pdb)}

	mux.HandleFunc("/", p.GetIndexHandler).Methods("GET")
	mux.HandleFunc("/project", p.GetProjectInfoHandler).Methods("GET")
	mux.HandleFunc("/artist", p.GetArtistInfoHandler).Methods("GET")
	mux.HandleFunc("/user", p.GetUserInfoHandler).Methods("GET")

	mux.HandleFunc("/auth/google/signup", p.googleLoginHandler)
	mux.HandleFunc("/auth/google/callback", p.googleAuthCallback)
	mux.HandleFunc("/auth/facebook/signup", p.facebookLoginHandler)
	mux.HandleFunc("/auth/facebook/callback", p.facebookAuthCallback)
	mux.HandleFunc("/auth/kakao/signup", p.kakaoLoginHandler)
	mux.HandleFunc("/auth/kakao/callback", p.kakaoAuthCallback)
	mux.HandleFunc("/auth/naver/signup", p.naverLoginHndler)
	mux.HandleFunc("/auth/naver/callback", p.naverAuthCallback)

	mux.HandleFunc("/user", p.PutUserInfoHandler).Methods("PUT")

	return mux

}
