package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	var reqpod ReqProjectsOfTheDay
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))

	reqpod.DemandDate = time.Now()
	reqpod.DemandPeriod = "1"

	data, err := json.Marshal(reqpod)
	if err != nil {
		log.Println("[ERR] json marshal err : ", err)
	}
	buff := bytes.NewBuffer(data)
	client := &http.Client{}

	req, err := http.NewRequest("GET", ts.URL+"/project/information", buff)

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	assert.Equal(res.StatusCode, http.StatusOK)

	respbytes, _ := ioutil.ReadAll(res.Body)
	log.Println("string resp body : ", string(respbytes))

	assert.Equal(http.StatusOK, res.StatusCode)

}

/*
func TestGetArtistInfoHandler(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))

	resp, err := http.Get(ts.URL + "/artist")
	if err != nil {
		log.Println("[ERR] http get err : ", err)
	}

	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestPutUserInfoHandler(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))

	var reqputuserinfo ReqJoinInfo
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for i := 1; i < 4; i++ {
		mediaheader := textproto.MIMEHeader{}
		mediaheader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"ID_photo_temp_%d.jpeg\"", i))
		mediaheader.Set("Content-ID", "media")
		mediaheader.Set("Content-Filename", fmt.Sprintf("ID_photo_temp_%d.jpeg", i))

		file, err := os.Open(fmt.Sprintf("../ID_photo_temp_%d.jpeg", i))
		if err != nil {
			log.Println("[ERR] os open err : ", err)
		}

		part, err := writer.CreatePart(mediaheader)
		if err != nil {
			log.Println("[ERR] create part err : ", err)
		}
		io.Copy(part, file)
	}

	reqputuserinfo.UserInfo.ID = "1123"
	reqputuserinfo.UserInfo.Name = "dukryung"
	reqputuserinfo.UserInfo.Nickname = "superunderdog"
	reqputuserinfo.UserInfo.Introduction = "dukryung is geniune"
	reqputuserinfo.UserInfo.AgreeEmailMarketing = true
	reqputuserinfo.AccountInfo.UserID = "777777"
	reqputuserinfo.AccountInfo.Bank = 1
	reqputuserinfo.AccountInfo.Account = "1111-111-1111-111111"
	reqputuserinfo.AccountInfo.AgreePolicy = true

	data, _ := json.Marshal(reqputuserinfo)

	metadataheader := textproto.MIMEHeader{}
	metadataheader.Set("Content-Type", "application/json")
	metadataheader.Set("Content-ID", "metadata")
	part, _ := writer.CreatePart(metadataheader)
	part.Write(data)

	writer.Close()
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/user", bytes.NewReader(body.Bytes()))
	if err != nil {
		log.Println("[ERR] new request err : ", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)

}

*/
