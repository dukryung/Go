package common

import (
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

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
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

func (p *programHandler) getProductHandler(w http.ResponseWriter, r *http.Request) {
	sessionid := getSessionId(r)
	productlist := p.db.getProducts(sessionid)
	rd.JSON(w, http.StatusOK, productlist)
}

func (p *programHandler) addProductHandler(w http.ResponseWriter, r *http.Request) {
	sessionid := getSessionId(r)
	name := r.FormValue("name")
	p.db.addProducts(name, sessionid)

}

func (p *programHandler) Close() {
	p.db.Close()
}

/*
func checkSign(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if strings.Contains(r.URL.Path, "sign") || strings.Contains(r.URL.Path, "auth") {
		next(w, r)
	}
}
*/
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

func MakeHandler(filepath string) *programHandler {
	r := mux.NewRouter()

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.NewStatic(http.Dir("public")))
	n.UseHandler(r)

	p := &programHandler{
		Handler: n,
		db:      NewDBHandler(filepath),
	}

	r.HandleFunc("/", p.isBasicHandler)
	r.HandleFunc("/", p.getProductHandler).Methods("GET")
	r.HandleFunc("/ws", wsHandler)

	return p

}
