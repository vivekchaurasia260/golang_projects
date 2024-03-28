package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("This is my First Program")

	var str string
	str = "Hello dear"
	fmt.Println(str)

	var b = true
	fmt.Println(b)

	d := 3.1415
	fmt.Println(d)

	var e int = int(d)
	fmt.Println(e)

	var f int8 = 43
	var g int16 = int16(f)
	fmt.Println(g)

	p, q := 7, 19
	r := p + q
	fmt.Println(r)

	// CLI APPLICATION
	fmt.Println("We are using fmt Package")
	in := bufio.NewReader(os.Stdin)
	s, _ := in.ReadString('\n')
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)
	fmt.Println(s + "!")

}
