package headServer

import (
	"fmt"
	"reflect"
)

type set struct {
	first  firstProperty
	second secondProperty
}

type firstProperty struct {
	third  thirdProperty
	fourth fourthProperty
}
type secondProperty struct{}

type thirdProperty struct{}
type fourthProperty struct{}

func ReflectElement() {

	s := &set{}
	element := reflect.ValueOf(*s)
	a := 4

	var x = interface{}(a)

	fmt.Println("element type: ", element.Type())

	for i := 0; i < element.NumField(); i++ {
		field := element.Field(i)
		fieldType := field.Type()
		fmt.Println("field type : ", fieldType)

		fmt.Println("field type name : ", fieldType.Name())
		value := reflect.New(fieldType)

		var methodFuncList []interface{}
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodFunc := value.MethodByName(method.Name).Interface()
			//f := reflect.ValueOf(methodFunc)
			//fmt.Println("num in :" ,f.Type().NumIn())
			//f.Call([]reflect.Value{})
			methodFuncList = append(methodFuncList, methodFunc)
		}

		xxx := reflect.ValueOf(x)
		yyy := []reflect.Value{}
		yyy = append(yyy,xxx)
		for k := 0; k < len(methodFuncList); k++ {
			f := reflect.ValueOf(methodFuncList[k])
			fmt.Println("num in :", f.Type().NumIn())

			//params := make([]reflect.Value, f.Type().NumIn())
			//fmt.Println("param length : ", len(params))
			f.Call(yyy)
		}

	}

}

func (f *firstProperty) PrintHello(a int) {
	fmt.Println("PrintHello a : ",a)
	fmt.Println("I am right!!!!")
}

func (f firstProperty) List(a int) {
	fmt.Println("List function a : ",a)
	fmt.Println("I am List function!!!!")
}