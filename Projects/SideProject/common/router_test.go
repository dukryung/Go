package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexPathHandler(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(MakeHandler("sideproject"))

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

}

func TestGetProjectInfoHandler(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))

	resp, err := http.Get(ts.URL + "/project?current_date=2021-03-22&demand_date=2021-03-22&demand_period=1")
	if err != nil {
		fmt.Println("@@@@@@", err)
	}

	respbytes, _ := ioutil.ReadAll(resp.Body)
	log.Println("string resp body : ", string(respbytes))

	assert.Equal(http.StatusOK, resp.StatusCode)

}
func TestGetArtistInfoHandler(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))

	resp, err := http.Get(ts.URL + "/artist")
	if err != nil {
		log.Println("@@@@@!!!!!!", err)
	}

	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestPutUserInfoHandler(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))

	var reqputuserinfo ReqJoinInfo
	client := http.Client{}

	reqputuserinfo.UserInfo.ID = "1123"
	reqputuserinfo.UserInfo.Name = "dukryung"
	reqputuserinfo.UserInfo.Nickname = "superunderdog"
	reqputuserinfo.UserInfo.Introduction = "dukryung is geniune"
	reqputuserinfo.UserInfo.AgreeEmailMarketing = true
	reqputuserinfo.AccountInfo.UserID = "777777"
	reqputuserinfo.AccountInfo.Bank = "1"
	reqputuserinfo.AccountInfo.Account = "1111-111-1111-111111"
	reqputuserinfo.AccountInfo.AgreePolicy = true

	data, _ := json.Marshal(reqputuserinfo)

	log.Println("!@#!@#!@#!@#!@#!@#!@#!@#!@#!@#!@#!@#:", string(data))
	b := bytes.NewBuffer(data)
	req, _ := http.NewRequest(http.MethodPut, ts.URL+"/user", b)

	res, _ := client.Do(req)
	assert.Equal(http.StatusOK, res.StatusCode)

}
