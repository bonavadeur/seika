package bonalib

import (
	"math/rand"
	"reflect"
	"strconv"
	"unsafe"

	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func Baka() string {
	return "Baka"
}

func RandNumber() string {
	return strconv.Itoa(rand.Intn(1000))
}

// 1;31: red
// 1;32: green
// 1;33: yellow
// 1;34: blue
// 1;35: purple

func Log(msg string, obj ...interface{}) {
	if msg == "" {
		msg = "-"
	}

	color := "\033[1;33m%v\033[0m" // yellow
	fmt.Printf("\033[1;33m0---bonaLog %v \033[0m", msg)

	for _, v := range obj {
		fmt.Printf("\033[1;33m%v \033[0m", v)
	}

	color = "\033[0m%v\033[0m" // reset
	fmt.Printf(color, "\n\n")
}

func Succ(msg string, obj ...interface{}) {
	if msg == "" {
		msg = "-"
	}

	color := "\033[1;32m%v\033[0m" // yellow
	fmt.Printf("\033[1;32m0---bonaLog %v \033[0m", msg)

	for _, v := range obj {
		fmt.Printf("\033[1;32m%v \033[0m", v)
	}

	color = "\033[0m%v\033[0m" // reset
	fmt.Printf(color, "\n\n")
}

func Warn(msg string, obj ...interface{}) {
	if msg == "" {
		msg = "-"
	}

	color := "\033[1;31m%v\033[0m" // yellow
	fmt.Printf("\033[1;31m0---bonaLog %v \033[0m", msg)

	for _, v := range obj {
		fmt.Printf("\033[1;31m%v \033[0m", v)
	}

	color = "\033[0m%v\033[0m" // reset
	fmt.Printf(color, "\n\n")
}

func Info(msg string, obj ...interface{}) {
	if msg == "" {
		msg = "-"
	}

	color := "\033[1;34m%v\033[0m" // yellow
	fmt.Printf("\033[1;34m0---bonaLog %v \033[0m", msg)

	for _, v := range obj {
		fmt.Printf("\033[1;34m%v \033[0m", v)
	}

	color = "\033[0m%v\033[0m" // reset
	fmt.Printf(color, "\n\n")
}

func Vio(msg string, obj ...interface{}) {
	if msg == "" {
		msg = "-"
	}

	color := "\033[1;35m%v\033[0m" // yellow
	fmt.Printf("\033[1;35m0---bonaLog %v \033[0m", msg)

	for _, v := range obj {
		fmt.Printf("\033[1;35m%v \033[0m", v)
	}

	color = "\033[0m%v\033[0m" // reset
	fmt.Printf(color, "\n\n")
}

func Line() {
	fmt.Printf("\n\n\n")
}

func Use(variable ...interface{}) {}

func Type(variable interface{}) string {
	return reflect.TypeOf(variable).String()
}

func Size(variable interface{}) int {
	size := unsafe.Sizeof(variable)
	return int(size)
}

func Logln(msg string, obj interface{}) {
	// 32: green 33: yellow
	if msg == "" {
		msg = "-"
	}
	if obj == "" {
		obj = "-"
	}
	color := "\033[1;33m%v\033[0m" // yellow
	str := spew.Sprintln("0---bonaLog", msg, obj)
	fmt.Printf(color, str)
	color = "\033[0m%v\033[0m" // reset
	fmt.Printf(color, "-------------------------------\n\n")
}
