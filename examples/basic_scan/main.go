package main

import (
	"fmt"
	"github.com/zkbupt/afrog"
)

func main() {
	if err := afrog.NewScanner([]string{"http://example.com"}, afrog.Scanner{}); err != nil {
		fmt.Println(err.Error())
	}
}
