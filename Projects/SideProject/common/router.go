package common

import (
	"Go/Projects/SideProject/common/database"
	"net/http"

	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

type project struct {
	db database.DBHandler
}

func (p *project) mainpagehandler(r http.ResponseWriter, w *http.Request) {

}

func MakeHandler() {
	r := mux.NewRouter()

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.NewStatic(http.Dir("public")))
	n.UseHandler(r)

	pdb := &database.ProjectDB{}

	p := &project{db: database.MakeDBHandler(pdb)}

	r.HandleFunc("/", p.mainpagehandler).Methods("GET")

}
