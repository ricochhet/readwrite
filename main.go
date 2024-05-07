package main

import (
	"fmt"

	"github.com/ricochhet/readwrite/readwrite"
)

func main() {
	fmt.Println(string(readwrite.Utf8ToUtf16("aaabbbccc")))
}
