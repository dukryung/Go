package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

func GenKey() (*rsa.PrivateKey, error) {
	keys, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("[ERR] failed to generate encrypt keys")
		return nil, err
	}

	return keys, nil
}

func Encrypt(text string, keys *rsa.PrivateKey) (string, error) {
	pk := &keys.PublicKey
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader,
		pk,
		[]byte(text))

	if err != nil {
		fmt.Println("[ERR] failed to encrypt text")
		return "", nil
	}

	return string(ciphertext), nil
}

func Decrypt(ciphertext string, privtekey *rsa.PrivateKey) (string, error) {
	plaintext, err := rsa.DecryptPKCS1v15(
		rand.Reader,
		privtekey,
		[]byte(ciphertext),
	)
	if err != nil {
		fmt.Println("[ERR] failed to decrypte text")
		return "", err
	}

	return string(plaintext), nil
}
