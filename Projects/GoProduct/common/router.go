package common

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/unrolled/render"

	"github.com/gorilla/sessions"

	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var rd *render.Render

type programHandler struct {
	http.Handler
	db DBHandler
}
type userHandler struct {
	http.Handler
	db DBHandler
}
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type SessionInfo struct {
	Issessionid bool `json:"issessionid"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func getSessionId(r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		panic(err)
	}

	val := session.Values["id"]

	return val.(string)
}

func (p *programHandler) isBasicHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/html/index.html", http.StatusTemporaryRedirect)
}
func (u *userHandler) readUserHandler(w http.ResponseWriter, r *http.Request) {
	sessionid := getSessionId(r)
	userlist := u.db.readProducts(sessionid)
	rd.JSON(w, http.StatusOK, userlist)
}

func (u *userHandler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userinfo *userinfo
	sessionid := getSessionId(r)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(userinfo)
	if err != nil {
		panic(err)
	}

	u.db.createUsers()

}

func (u *userHandler) updateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (u *userHandler) deleteUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (p *programHandler) readProductHandler(w http.ResponseWriter, r *http.Request) {
	sessionid := getSessionId(r)
	productlist := p.db.readProducts(sessionid)
	rd.JSON(w, http.StatusOK, productlist)
}

func (p *programHandler) createProductHandler(w http.ResponseWriter, r *http.Request) {
	sessionid := getSessionId(r)
	name := r.FormValue("name")
	p.db.createProducts(name, sessionid)

}

func (p *programHandler) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	sessionid := getSessionId(r)
	name := r.FormValue("name")
	p.db.updateProducts(name, sessionid)
}

func (p *programHandler) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	sessionid := getSessionId(r)
	name := r.FormValue("name")
	p.db.deleteProducts(name, sessionid)
}

func (p *programHandler) Close() {
	p.db.Close()
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	m := &Message{}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {

		err = conn.ReadJSON(m)
		if err != nil {
			log.Println(err)
			return
		}

	}
}

func checksign(w http.ResponseWriter, r *http.Request) bool {
	var issessionid bool
	sessionid := getSessionId(r)
	if sessionid != "" {
		issessionid = true
	} else {
		issessionid = false
	}
	return issessionid
}

func MakeHandler(dbfilepath DBfilepath) *programHandler {
	r := mux.NewRouter()

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.NewStatic(http.Dir("public")))
	n.UseHandler(r)

	pdb := &ProductDB{filepath: dbfilepath.Productdbfilepath}
	udb := &UserDB{filepath: dbfilepath.Userdbfilepath}

	p := &programHandler{
		Handler: n,
		db:      NewDBHandler(pdb),
	}

	u := &userHandler{
		Handler: n,
		db:      NewDBHandler(udb),
	}

	r.HandleFunc("/user", u.readUserHandler).Methods("GET")
	r.HandleFunc("/user", u.createUserHandler).Methods("POST")
	r.HandleFunc("/user", u.updateUserHandler).Methods("PUT")
	r.HandleFunc("/user", u.deleteUserHandler).Methods("DELETE")

	r.HandleFunc("/", p.isBasicHandler)
	r.HandleFunc("/products", p.readProductHandler).Methods("GET")
	r.HandleFunc("/products", p.createProductHandler).Methods("POST")
	r.HandleFunc("/products", p.updateProductHandler).Methods("PUT")
	r.HandleFunc("/products", p.deleteProductHandler).Methods("DELETE")
	r.HandleFunc("/auth/google/login", googleLoginHandler)
	r.HandleFunc("/ws", wsHandler)

	return p

}
