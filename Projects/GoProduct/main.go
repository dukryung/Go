package main

import (
	"Go/Projects/GoProduct/common"
	"net/http"
)

func main() {
	// read dbfilepath

	// read dbfilepath
	p := common.MakeHandler(common.DBfilepath)
	defer p.Close()

	http.ListenAndServe(":8080", p)
}
