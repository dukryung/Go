package main

import (
	"Go/GoWeb/GoWeb_3rd/myapp"
	"net/http"
)

func main() {

	http.ListenAndServe(":8080", myapp.NewHandler())
}
