package main

import (
	"Go/GoWeb/GoWeb_9th/app"
	"net/http"
)

func main() {
	m := app.MakeNewHandler("./test.db")
	defer m.Close()

	http.ListenAndServe(":8080", m)

}
