package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexHandler(t *testing.T) {

	assert := assert.New(t)
	client := http.Client{}

	ts := httptest.NewServer(MakeHandler())

	var userinfo info

	userinfo.Name = "hello dukryung"
	userinfo.Age = 15

	data, err := json.Marshal(userinfo)
	if err != nil {
		log.Println("erradsfadfadfasd")
	}

	b := bytes.NewBuffer(data)
	req, err := http.NewRequest("GET", ts.URL, b)
	if err != nil {
		log.Fatal("[request err : ", err)
	}

	res, _ := client.Do(req)
	assert.Equal(http.StatusOK, res.StatusCode)

	/*
		fd, err := os.Open("./test_1.txt")
		if err != nil {
			log.Fatal("os Open err : ", err)
		}
		defer fd.Close()

		fd2, err := os.Open("./test_2.txt")
		if err != nil {
			log.Fatal("os Open err : ", err)
		}
		defer fd2.Close()

		b := bytes.Buffer{}
		writer := multipart.NewWriter(&b)

		part, err := writer.CreateFormFile("test", "test_1.txt")
		_, err = io.Copy(part, fd)
		if err != nil {
			log.Fatal("io Copy err : ", err)
		}
		err = writer.WriteField("file", "test_file")
		if err != nil {
			log.Fatal("WriteField err : ", err)
		}

		part, err = writer.CreateFormFile("test", "test_2.txt")
		io.Copy(part, fd2)

		writer.Close()

		log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! ts URL : ", ts.URL)
		req, err := http.NewRequest("POST", ts.URL+"/test", &b)
		if err != nil {
			log.Fatal("[request err : ", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}

		res, err := client.Do(req)
		if err != nil {
			log.Fatal("client Do err : ", err)
		}

	*/
	/*
		res, err := http.Get(ts.URL + "/")
	*/
	/*
		assert.NoError(err)
		assert.Equal(http.StatusOK, res.StatusCode)
	*/
}
