package main

import (
	"Go/Projects/GoProduct/common"
	"net/http"
)

func main() {
	p := common.MakeHandler("./info.db")
	defer p.Close()

	http.ListenAndServe(":8080", p)
}
