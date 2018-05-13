package main

import (
	"fmt"
	"os"
)

func main() {
	len(os.Args)

	fmt.Println("vim-go")
}

type Counter struct {
	Sentence string `json:"sencence"`
	Count    int    `json:"count"`
}
