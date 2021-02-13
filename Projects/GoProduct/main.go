package main

import (
	"Go/Projects/GoProduct/common"
	"net/http"
)

var ConfigFileDir = "./Config"
var ConfigFileName = "treasure_house_encrypt_filename"

func main() {
	// read dbfilepath

	EncryptFileName := common.EncryptConfigFileName(ConfigFileName)
	f := &common.FileInfo{FilePath: ConfigFileDir, FileName: EncryptFileName}
	dbfilepath := f.CreateConfigFile()
	// read dbfilepath
	p := common.MakeHandler(*dbfilepath)
	defer p.Close()

	http.ListenAndServe(":8080", p)
}
