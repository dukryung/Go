package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

func Encrypt(text string) (string, error) {
	keys, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("[ERR] failed to generate encrypt keys")
		return "", nil
	}

	publickey := keys.PublicKey
}

func Decrypt(text string) (string, error) {

}
