package controlfile

import (
	"fmt"
	"os"
)

type fileinfo struct {
	Filepath string
	FileName string
}

type folderinfo struct {
	Foldername string
}

func createfile(fileinfo *fileinfo) {
	isfile := checkFileExist(fileinfo)
	if isfile == "exist" {
		return
	}

	fd, err := os.Create(fileinfo.Filepath + fileinfo.FileName)
	if err != nil {
		fmt.Println("[LOG] Create file Error :", err)
		return
	}
	defer fd.Close()
}

func createfolder(folderinfo *folderinfo) {
	isfolder := checkFolderExist(folderinfo)
	if isfolder == "exist" {
		return
	}

	err := os.MkdirAll(folderinfo.Foldername, 0777)
	if err != nil {
		fmt.Println("[LOG] Create folder Error :", err)
		return
	}
}

func checkFileExist(fileinfo *fileinfo) string {
	if _, err := os.Stat(fileinfo.Filepath + fileinfo.FileName); os.IsNotExist(err) {
		return "notexist"
	} else {
		return "exist"
	}
}

func checkFolderExist(folderinfo *folderinfo) string {
	if _, err := os.Stat(folderinfo.Foldername); os.IsNotExist(err) {
		return "notexist"
	} else {
		return "exist"
	}
}

func writefile(fileinfo *fileinfo, text string) {

	fd, err := os.Create(fileinfo.Filepath + fileinfo.FileName)
	if err != nil {
		fmt.Println("[LOG] Create file Error :", err)
		return
	}

	defer fd.Close()

	_, err = fd.WriteString(text)
	if err != nil {
		fmt.Println("[LOG] WriteString err:", err)
		return
	}
	fd.Sync()
}
