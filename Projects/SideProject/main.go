package main

import (
	"Go/Projects/SideProject/common"
	"net/http"

	"github.com/urfave/negroni"
)

func main() {

	databasename := "sideproject"

	mux := common.MakeHandler(databasename)
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.NewStatic(http.Dir("public")))
	n.UseHandler(mux)

	err := http.ListenAndServe(":8080", n)
	if err != nil {
		panic(err)
	}

}
