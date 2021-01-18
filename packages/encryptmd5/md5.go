package encryptmd5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func Encrypt(text string) string {
	data := []byte(text)
	fmt.Println("[LOG] md5 encrypt text :", md5.Sum(data))
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
