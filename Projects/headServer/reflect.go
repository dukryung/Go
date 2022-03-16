package headServer

import (
	"fmt"
	"reflect"
	"strings"
)

type set struct {
	first  firstProperty
	second secondProperty
}

type firstProperty struct {
}
type secondProperty struct {
}

func(f firstProperty) List(){

}


func ReflectElement() {

	s := &set{}
	element := reflect.ValueOf(*s)

	fmt.Println("element type: ",element.Type())

	for i := 0; i <element.NumField(); i++ {
		field  := element.Field(i)
		fieldType := field.Type()
		fmt.Println("f type : ",fieldType)


		fieldName := strings.ToLower(fieldType.Name())
		fmt.Println("f type name : ",fieldName)

		value := reflect.New(fieldType)
		for j := 0 ; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			fmt.Println("method name: ",method.Name)
		}

	}

}