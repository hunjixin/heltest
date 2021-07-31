package main

import (
	"fmt"
	"strings"
)

func main() {
	xxx := "a123.2131_123123.sdverjkcsdfasd___123123_11111111"
	vv := strings.Split(xxx, "___")
	fmt.Println(vv)
}
