package main

import (
	"Go/packages/encryptmd5"
	"fmt"

	"github.com/google/uuid"
)

var encryptkeyfilename = "treasurehouse_admin_keyfile"

func main() {
	filename := encryptmd5.Encrypt(encryptkeyfilename)
	fmt.Println("[LOG] filename:", filename)

	id := uuid.New()
	fmt.Println(id.String())
}
