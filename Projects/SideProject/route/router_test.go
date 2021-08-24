package route

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"sideproject/route/iamport"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestUser struct {
	ID int `json:"user_id"`
}
type TestArist struct {
	ArtistID int `json:"artist_id"`
}

/*
func TestIndexPathHandler(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(MakeHandler("sideproject"))

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

}

func TestGetProjectInfoHandler(t *testing.T) {
	var reqpod = &ReqProjectsOfTheDay{}
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))

	now := time.Now()
	reqpod.DemandDate = now

	reqpod.DemandPeriod = 1

	data, err := json.Marshal(reqpod)
	if err != nil {
		log.Println("[ERR] json marshal err : ", err)
	}
	buff := bytes.NewBuffer(data)
	client := &http.Client{}

	req, err := http.NewRequest("GET", ts.URL+"/project/information", buff)

	//req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	respbytes, _ := ioutil.ReadAll(res.Body)
	log.Println("project list : ", string(respbytes))

	assert.Equal(http.StatusOK, res.StatusCode)

}

func TestGetArtistInfoHandler(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))

	res, err := http.Get(ts.URL + "/artist")
	if err != nil {
		log.Println("[ERR] http get err : ", err)
	}
	respbytes, _ := ioutil.ReadAll(res.Body)
	log.Println("artist list : ", string(respbytes))

	assert.Equal(http.StatusOK, res.StatusCode)
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
		mediaheader.Set("Content-ID", "userimage")
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

	reqputuserinfo.UserInfo.ID = 1
	reqputuserinfo.UserInfo.Name = "dukryung_1"
	reqputuserinfo.UserInfo.Nickname = "duck_1"
	reqputuserinfo.UserInfo.Email = "dukryung_2@naver.com"
	reqputuserinfo.UserInfo.Introduction = "introduction_1"
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



func TestGetProfileFramInfoHandler(t *testing.T) {
	var userinfo TestUser
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	userinfo.ID = 111

	data, err := json.Marshal(userinfo)
	if err != nil {
		log.Println("[ERR] json marshal err")
	}

	body := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/profileuser/frame", body)
	if err != nil {
		log.Println("[ERR] new request err")
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client err :", err)
	}
	assert.Equal(http.StatusOK, res.StatusCode)
	frameinfo, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("[ERR] read all data err : ", err)
	}

	log.Println("[LOG] response data :", string(frameinfo))
}

func TestGetProfileProjectHandler(t *testing.T) {
	var userinfo TestUser
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	userinfo.ID = 111

	data, err := json.Marshal(userinfo)
	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/profileuser/project", buff)
	if err != nil {
		log.Println("[ERR] new request err : ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	projectinfo, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("[ERR] read all err : ", err)
	}

	log.Println("[LOG] project info : ", string(projectinfo))

	assert.Equal(http.StatusOK, res.StatusCode)

}

func TestGetProfileSellHandler(t *testing.T) {
	var userinfo TestUser

	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	userinfo.ID = 111

	data, err := json.Marshal(userinfo)
	if err != nil {
		log.Println("[ERR] json marshal err : ", err)
	}

	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/profileuser/sell", buff)
	if err != nil {
		log.Println("[ERR] new requesterr : ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)

	sellhistoryinfo, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("[ERR] read all err : ", err)
	}

	log.Println("[LOG] sell history info : ", string(sellhistoryinfo))
}

func TestGetProfileBuyHandler(t *testing.T) {
	var userinfo TestUser
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	userinfo.ID = 111

	data, err := json.Marshal(userinfo)
	if err != nil {
		log.Println("[ERR] json marshal err : ", err)
	}

	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/profileuser/buy", buff)
	if err != nil {
		log.Println("[ERR] request err : ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client Do err : ", err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)

	buyhistoryinfo, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("[ERR] read all err : ", err)

	}

	log.Println("[LOG] buy history info err : ", string(buyhistoryinfo))

}
func TestGetProfileWithdrawHandler(t *testing.T) {
	var userinfo = *&TestUser{}

	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))
	client := http.Client{
		Timeout: time.Second * 10,
	}

	userinfo.ID = 111

	data, err := json.Marshal(userinfo)
	if err != nil {
		log.Println("[ERR] json marshal err: ", err)
	}

	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/profileuser/withdraw", buff)
	if err != nil {
		log.Println("[ERR] new request err : ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)
	withdrawhistoryinfo, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("[ERR] read all err : ", err)
	}

	log.Println("[LOG] withdraw history : ", string(withdrawhistoryinfo))

}


func TestReadModificationUserInfo(t *testing.T) {
	var userinfo = &TestUser{}

	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))
	client := http.Client{
		Timeout: time.Second * 10,
	}

	userinfo.ID = 111

	data, err := json.Marshal(userinfo)
	if err != nil {
		log.Println("[ERR] json marshal err: ", err)
	}

	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/profileuser/modification", buff)
	if err != nil {
		log.Println("[ERR] new request err : ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)
	withdrawhistoryinfo, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("[ERR] read all err : ", err)
	}

	log.Println("[LOG] withdraw history : ", string(withdrawhistoryinfo))

}

func TestPutModificationUserInfoHandler(t *testing.T) {
	var reqmodinfo = &ReqModificationUserInfo{}

	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))
	client := http.Client{
		Timeout: time.Second * 10,
	}

	reqmodinfo.ID = 111
	reqmodinfo.Name = "testname"
	reqmodinfo.Nickname = "testnickname"
	reqmodinfo.Email = "test@naver.com"
	reqmodinfo.AgreeEmailMarketing = true
	reqmodinfo.Introduction = "testintroduction"
	reqmodinfo.ImageLink = "test/image"

	data, err := json.Marshal(reqmodinfo)
	if err != nil {
		log.Println("[ERR] json marshal err: ", err)
	}

	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodPut, ts.URL+"/profileuser/modification", buff)
	if err != nil {
		log.Println("[ERR] new request err : ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)

}

func TestGetProfileArtistInfoHandler(t *testing.T) {
	var artistinfo = &TestArist{}

	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))
	client := http.Client{
		Timeout: time.Second * 10,
	}

	artistinfo.ArtistID = 111

	data, err := json.Marshal(artistinfo)
	if err != nil {
		log.Println("[ERR] json marshal err: ", err)
	}

	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/profileartist/information", buff)
	if err != nil {
		log.Println("[ERR] new request err : ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)
	artistprofileinfo, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("[ERR] read all err : ", err)
	}

	log.Println("[LOG] artistprofileinfo history : ", string(artistprofileinfo))

}

func TestGetPersonalInformationHandler(t *testing.T) {
	var userinfo = &TestUser{}

	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))
	client := http.Client{
		Timeout: time.Second * 10,
	}

	userinfo.ID = 111

	data, err := json.Marshal(userinfo)
	if err != nil {
		log.Println("[ERR] json marshal err: ", err)
	}

	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/personal/information", buff)
	if err != nil {
		log.Println("[ERR] new request err : ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)
	personalinfo, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("[ERR] read all err : ", err)
	}

	log.Println("[LOG] personalinfo history : ", string(personalinfo))

}
func TestPutPersonalInformationHandler(t *testing.T) {
	var reqpersonalinformation = &ReqPersonalInformation{}

	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))
	client := http.Client{
		Timeout: time.Second * 10,
	}

	reqpersonalinformation.Account.UserID = 111
	reqpersonalinformation.Account.Bank = 1
	reqpersonalinformation.Account.Account = "1111-111-1111-11111"

	data, err := json.Marshal(reqpersonalinformation)
	if err != nil {
		log.Println("[ERR] json marshal err: ", err)
	}

	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodPut, ts.URL+"/personal/information", buff)
	if err != nil {
		log.Println("[ERR] new request err : ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)
}
*/

/*
func TestGetProjectDetailProjectInformationHandler(t *testing.T) {
	var reqprojectdetailinfo = &ReqProjectDetailInfo{}
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("sideproject"))
	client := http.Client{
		Timeout: time.Second * 10,
	}

	reqprojectdetailinfo.UserID = 111
	reqprojectdetailinfo.ProjectID = 551

	data, err := json.Marshal(reqprojectdetailinfo)
	if err != nil {
		log.Println("[ERR] json marshal err: ", err)
	}

	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/project/information/detail/project", buff)
	if err != nil {
		log.Println("[ERR] new request err : ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)
	resprojectdetailinfo, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("[ERR] read all err : ", err)
	}

	log.Println("[LOG] personalinfo history : ", string(resprojectdetailinfo))

}

*/
/*
func TestParentiamport(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(MakeHandler("sideproject"))
	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/iamport", nil)
	if err != nil {
		log.Println("[ERR] new request err : ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)

	token, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("[ERR] read all err : ", err)
	}

	log.Println("[LOG] token history : ", string(token))
}

*/
func TestCheckPAYMT(t *testing.T) {
	reqsrc := &iamport.REQRSC{Iid: "imp_367251774536",
		Mid: "xxx"}
	assert := assert.New(t)

	ts := httptest.NewServer(MakeHandler("sideproject"))
	client := http.Client{
		Timeout: time.Second * 10,
	}

	data, err := json.Marshal(reqsrc)
	if err != nil {
		log.Println("[ERR] json marshal err: ", err)
	}

	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/iamport/payment/complete", buff)
	if err != nil {
		log.Println("[ERR] new request err : ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[ERR] client do err : ", err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)

}
