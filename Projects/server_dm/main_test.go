package main

import (
	"testing"
	"unicode"
)

/*
func TestGetFriends(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(Makeroutehandler())
	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL+"/friends", nil)
	if err != nil {
		panic(err)
	}

	xpubkey := "1234"
	encodebase64 := base64.StdEncoding.EncodeToString([]byte(xpubkey))
	req.Header.Set("X-PubKey", encodebase64)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	t.Log(string(data))
	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestGetFriendsMatchedTag(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(Makeroutehandler())
	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL+"/friends/tag/school", nil)
	if err != nil {
		panic(err)
	}

	xpubkey := "1234"
	encodebase64 := base64.StdEncoding.EncodeToString([]byte(xpubkey))
	req.Header.Set("X-PubKey", encodebase64)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	log.Println(string(data))
	t.Log(string(data))
	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestGetTags(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(Makeroutehandler())
	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL+"/friends/tags", nil)
	if err != nil {
		panic(err)
	}

	xpubkey := "1234"
	encodebase64 := base64.StdEncoding.EncodeToString([]byte(xpubkey))
	req.Header.Set("X-PubKey", encodebase64)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	log.Println(string(data))
	t.Log(string(data))
	assert.Equal(http.StatusOK, res.StatusCode)
}
*/
/*
func TestDelTagsMatchedFriendTag(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(Makeroutehandler())
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", ts.URL+"/friends/tags/work", nil)
	if err != nil {
		panic(err)
	}

	xpubkey := "1234"
	encodebase64 := base64.StdEncoding.EncodeToString([]byte(xpubkey))
	req.Header.Set("X-PubKey", encodebase64)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)
}
*/
/*
func TestPutTags(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(Makeroutehandler())
	client := &http.Client{}

	taginfo := struct {
		Name string `json:"name"`
	}{
		Name: "newtag",
	}

	data, err := json.Marshal(taginfo)
	if err != nil {
		panic(err)
	}
	buff := bytes.NewBuffer(data)
	req, err := http.NewRequest("PUT", ts.URL+"/friends/tags/work", buff)
	if err != nil {
		panic(err)
	}

	xpubkey := "1234"
	encodebase64 := base64.StdEncoding.EncodeToString([]byte(xpubkey))
	req.Header.Set("X-PubKey", encodebase64)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)
}
*/
/*
func TestPostFriend(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(Makeroutehandler())
	client := &http.Client{}

	friendinfo := struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}{
		Name: "newtag",
		Tags: []string{"grandparent", "korea"},
	}

	data, err := json.Marshal(friendinfo)
	if err != nil {
		panic(err)
	}
	buff := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", ts.URL+"/friends/:4", buff)
	if err != nil {
		panic(err)
	}

	xpubkey := "1234"
	encodebase64 := base64.StdEncoding.EncodeToString([]byte(xpubkey))
	req.Header.Set("X-PubKey", encodebase64)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)
}
func TestPutFriend(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(Makeroutehandler())
	client := &http.Client{}

	friendinfo := struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
		}{
			Name: "newtag",
			Tags: []string{"grandparent", "korea"},
		}

		data, err := json.Marshal(friendinfo)
		if err != nil {
			panic(err)
		}
		buff := bytes.NewBuffer(data)
		req, err := http.NewRequest("PUT", ts.URL+"/friends/1", buff)
		if err != nil {
			panic(err)
		}

		xpubkey := "1234"
		encodebase64 := base64.StdEncoding.EncodeToString([]byte(xpubkey))
		req.Header.Set("X-PubKey", encodebase64)
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		assert.Equal(http.StatusOK, res.StatusCode)
	}
*/
/*
func TestDelFriend(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(Makeroutehandler())
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", ts.URL+"/friends/1", nil)
	if err != nil {
		panic(err)
	}

	xpubkey := "1234"
	encodebase64 := base64.StdEncoding.EncodeToString([]byte(xpubkey))
	req.Header.Set("X-PubKey", encodebase64)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestDelFriends(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(Makeroutehandler())
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", ts.URL+"/friends", nil)
	if err != nil {
		panic(err)
	}

	xpubkey := "1234"
	encodebase64 := base64.StdEncoding.EncodeToString([]byte(xpubkey))
	req.Header.Set("X-PubKey", encodebase64)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	assert.Equal(http.StatusOK, res.StatusCode)
}
*/
/*
type User struct {
	Name string `validate:"contains=@"`
	Age  string
}

func TestValid(t *testing.T) {

	validate := validator.New()
	validate.
	var u = User{}
	u.Name = "@dukryung"
	err := validate.Struct(u)
	if err != nil {
		t.Log(err)
	}

	t.Log(u)
}
*/

func TestTest(t *testing.T) {

	var id = "drc0830_.>"

	for _, char := range id {
		if !unicode.IsDigit(char) && !unicode.IsLetter(char) && char != '_' && char != '.' {
			t.Log("wrong regex : ", string(char))
		}

	}

}
