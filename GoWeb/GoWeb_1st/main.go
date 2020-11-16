package main

import (
	"GoWeb/myapp"
	"net/http"
)

func main() {

	http.ListenAndServe(":8080", myapp.NewHttpHandler())

}
