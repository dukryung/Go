package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"

	"google.golang.org/api/option"

	firebase "firebase.google.com/go/v4"
)

//func main() {

// config := &firebase.Config{
// 	ProjectID:     "sideproject-308610",
// 	StorageBucket: "sideproject-308610.appspot.com",
// }
// opt := option.WithCredentialsFile("./sideproject-308610-firebase-adminsdk-wmpyn-8b162af899.json")
// app, err := firebase.NewApp(context.Background(), config, opt)
// if err != nil {
// 	log.Fatal(err)
// }

// client, err := app.Storage(context.Background())
// if err != nil {
// 	log.Fatal("-------------2", err)
// }

// // fd, err := os.Open("./test.txt")
// // if err != nil {
// // 	log.Fatal("------------2", err)
// // }
// // defer fd.Close()

// bucket, err := client.DefaultBucket()
// if err != nil {
// 	log.Fatal("------------3", err)
// }

// src1 := bucket.Object("test.txt")

// dst1 := bucket.Object("user")
// attr, err := dst1.ComposerFrom(src1).Run(context.Background())
// if err != nil {
// 	log.Fatal("------------3", err)
// }

//wc := bucket.Object("dfasdfasdf").NewWriter(context.Background())

//wc := bucket.Object("user/test.txt").NewWriter(context.Background())
//_, err = io.Copy(wc, fd)
//if err != nil {
//	log.Fatal("---------4", err)
//}

//defer wc.Close()

// for i := 1; i < 6; i++ {
// 	fileaname := fmt.Sprintf("test_%s.txt", i)
// 	fd, err := os.Open("./" + fileaname)
// 	if err != nil {
// 		log.Fatal("[ERR] os Open", err)
// 	}
// 	defer fd.Close()

// 	body := &bytes.Buffer{}

// 	writer := multipart.NewWriter(body)
// 	part, err := writer.CreateFormFile("my_file", fileaname)
// 	if err != nil {
// 		log.Fatal("[ERR] CreateFormFile :", err)
// 	}

// 	io.Copy(part, fd)
// }

// config := &firebase.Config{
// 	ProjectID:     "sideproject-308610",
// 	StorageBucket: "sideproject-308610.appspot.com",
// }

// opt := option.WithCredentialsFile("./sideproject-308610-firebase-adminsdk-wmpyn-8b162af899.json")
// app, err := firebase.NewApp(context.Background(), config, opt)
// if err != nil {
// 	log.Fatal(err)
// }

// client, err := app.Storage(context.Background())
// if err != nil {
// 	log.Fatal("-------------2", err)
// }

// bucket, err := client.DefaultBucket()
// if err != nil {
// 	log.Fatal("------------3", err)
// }

// var wc *storage.Writer
// for i := 1; i < 6; i++ {
// 	data, err := ioutil.ReadFile(fmt.Sprintf("./test_%d.txt", i))
// 	if err != nil {
// 		log.Fatal("os Open err : ", err)
// 	}

// 	wc = bucket.Object(fmt.Sprintf("./test_%d.txt", i)).NewWriter(context.Background())
// 	wc.Write(data)

// }

//time.Sleep(time.Minute * 10)
/*
	write := multipart.NewWriter()
	mux := mux.NewRouter()

	mux.HandleFunc("/", indexhandler)
	mux.HandleFunc("/upload", getinfohandler).Methods("POST")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
*/
//}

/*
func indexhandler(w http.ResponseWriter, r *http.Request) {
	log.Println("!!!!!!!!!!!!!!!!!!!!")
	var b bytes.Buffer

	x := multipart.NewWriter(&b)

	value := map[string]io.Reader{
		"file": mustOpen("test_1.txt"),
	}

	for key, t := range value {

		fw, err := x.CreateFormFile(key, t.(*os.File).Name())
		if err != nil {
			log.Fatal("createform err : ", err)
		}

		_, err = io.Copy(fw, t)
		if err != nil {
			log.Fatal("err : ", err)
		}
	}
	log.Println("!!!!!!!!!!!!!!!!!!!!2")
	req, err := http.NewRequest("POST", "/upload", &b)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := &http.Client{}
	req.Header.Set("Content-Type", x.FormDataContentType())
	log.Println("!!!!!!!!!!!!!!!!!!!!3")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("client do err", err)
	}

	log.Println("!!!!!!!!!!!!!:", res.StatusCode)

}

func getinfohandler(w http.ResponseWriter, r *http.Request) {
	log.Println("r.body : ", r.Body)

	fmt.Fprintf(w, "hellworld")
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func multipart() {
	file, err := os.Open("test_1.txt")
	if err != nil {
		log.Fatal("os.Open err : ", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
*/
/*
var rd *render.Render = render.New()
*/
func main() {

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

	config := &firebase.Config{
		ProjectID:     "sideproject-308610",
		StorageBucket: "sideproject-308610.appspot.com",
	}
	opt := option.WithCredentialsFile("./sideproject-308610-firebase-adminsdk-wmpyn-8b162af899.json")
	app, _ := firebase.NewApp(context.Background(), config, opt)

	client, _ := app.Storage(context.Background())
	bucket, _ := client.DefaultBucket()

	wc := bucket.Object("user/").NewWriter(context.Background())

	_, err = io.Copy(wc, &b)
	if err != nil {
		log.Fatal("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@:", err)
	}
	wc.Close()

	/*
		err := http.ListenAndServe(":8080", MakeHandler())
		if err != nil {
			log.Fatal("listenandserve err : ", err)
		}
	*/
}

// func test(r *http.Request) {
// 	_, fileheader, _ := r.FormFile("asdf")

// 	for _, fd := range r.MultipartForm.File["sdfsdf"] {
// 		fd.Open()
// 	}
// }

/*
func MakeHandler() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/test", IndexHandler).Methods("POST")

	return mux

}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
*/
/*
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
*/

//rd.JSON(w, http.StatusOK, nil)
/*
	log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!")

	//log.Println("test test test : ", r.MultipartForm.Value["file"])
	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, h := range r.MultipartForm.File["test"] {

		filename := h.Filename
		log.Println("!@#!@#!@#!@#!@#!@#!@#!@#!@#!@#!@#:", filename)
		file, err := h.Open()
		if err != nil {
			log.Fatal(err)
		}
		file.Close()

	}
}
*/
