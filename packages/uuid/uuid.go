package uuid

import (
	"fmt"

	"github.com/google/uuid"
)

func Getuuid() string {

	id := uuid.New()
	fmt.Println(id.String())
	return id.String()
}
