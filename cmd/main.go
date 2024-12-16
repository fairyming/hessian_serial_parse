package main

import (
	"flag"
	"fmt"
	"os"

	hsParse "github.com/fairyming/hessian_serial_parse"
	"github.com/fairyming/hessian_serial_parse/utils"
)

var Path = flag.String("path", "", "path")

func main() {

	flag.Parse()
	data, err := os.ReadFile(*Path)
	if err != nil {
		panic(err)
	}

	parse, err := hsParse.NewHessianParse(data)
	if err != nil {
		panic(err)
	}
	result, err := parse.Parse()
	if err != nil {
		panic(err)
	}

	jsonData, err := utils.EncodeJson(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonData)
}
