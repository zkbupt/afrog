package main

import (
	"fmt"
	"github.com/zkbupt/afrog"
)

func main() {
	if err := afrog.NewScanner([]string{}, afrog.Scanner{
		TargetsFile: "./urls.txt",
	}); err != nil {
		fmt.Println(err.Error())
	}
}
