package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"

	"golang.org/x/crypto/pbkdf2"
)

type FileInfo struct {
	FilePath string
	FileName string
}

func EncryptConfigFileName(keyfilename string) string {

	data := []byte(keyfilename)
	hash := md5.Sum(data)
	filename := hex.EncodeToString(hash[:])

	return filename
}

func (f *FileInfo) CreateConfigFile() *DBfilepath {
	var fd *os.File
	var err error

	var encryptedcfginfo string
	var decryptedcfginfo string
	var dbfilepath DBfilepath

	isfolder := f.CheckFolderExist()
	if isfolder != "exist" {
		f.CreateConfigFolder()
	}

	isfile := f.CheckFileExist()
	if isfile != "exist" {
		fd, err = os.Create(f.FilePath + "/" + f.FileName)
		if err != nil {
			fmt.Println("[LOG] Create file Error :", err)
			return nil
		}

		origincfginfo := f.CreateOriginConfigInfo()
		encryptedcfginfo = f.EcryptConfigInfo(origincfginfo)
		f.WriteConfigFile(encryptedcfginfo, fd)
	} else {
		data, err := ioutil.ReadFile(f.FilePath + "/" + f.FileName)
		if err != nil {
			fmt.Println("[LOG] ReadFile err :", err)
			return nil
		}
		encryptedcfginfo = string(data)
	}

	decryptedcfginfo = f.DecryptConfigInfo(encryptedcfginfo)

	//{-----------read toml config information-----------------
	_, err = toml.Decode(decryptedcfginfo, &dbfilepath)
	if err != nil {
		fmt.Println("[LOG] decode toml err :", err)
		return nil
	}
	//------------read toml config information-----------------}

	defer fd.Close()

	return &dbfilepath
}

func (f *FileInfo) CreateConfigFolder() {

	err := os.MkdirAll(f.FilePath, 0777)
	if err != nil {
		fmt.Println("[LOG] Create folder Error :", err)
		return
	}
}

func (f *FileInfo) CheckFileExist() string {
	if _, err := os.Stat(f.FilePath + "/" + f.FileName); os.IsNotExist(err) {
		return "notexist"
	} else {
		return "exist"
	}
}

func (f *FileInfo) CheckFolderExist() string {
	if _, err := os.Stat(f.FilePath); os.IsNotExist(err) {
		return "notexist"
	} else {
		return "exist"
	}
}

func (f *FileInfo) WriteConfigFile(text string, fd *os.File) {
	var err error

	fd, err = os.Create(f.FilePath + "/" + f.FileName)
	if err != nil {
		fmt.Println("[LOG] Create file Error (c *ConfigInfo) :", err)
		return
	}

	_, err = fd.WriteString(text)
	if err != nil {
		fmt.Println("[LOG] WriteString err:", err)
		return
	}
	fd.Sync()
}

func (f *FileInfo) CreateOriginConfigInfo() string {
	var tomlinfo = `Productdbfilepath = "./db/product.db"
					Userdbfilepath = "./db/user.db"`
	return tomlinfo
}

func (f *FileInfo) DeriveKey() []byte {

	password := "house"
	salt := make([]byte, 8)
	salt = []byte("treasure")

	return pbkdf2.Key([]byte(password), salt, 1000, 32, sha256.New)
}

func (f *FileInfo) EcryptConfigInfo(cfginfo string) string {

	key := f.DeriveKey()
	iv := make([]byte, 12)
	iv = []byte("treasure_hou")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aesgcm, err := cipher.NewGCM(block)
	data := aesgcm.Seal(nil, iv, []byte(cfginfo), nil)

	return hex.EncodeToString(data)

}

func (f *FileInfo) DecryptConfigInfo(encryptedcfginfo string) string {
	encrypteddata, err := hex.DecodeString(encryptedcfginfo)
	if err != nil {
		panic(err)
	}

	key := f.DeriveKey()
	iv := make([]byte, 12)
	iv = []byte("treasure_hou")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	decrypteddata, err := aesgcm.Open(nil, iv, encrypteddata, nil)
	if err != nil {
		panic(err)
	}
	return string(decrypteddata)
}
